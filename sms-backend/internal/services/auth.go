package services

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"sms-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthClaims struct {
	UserID    uint   `json:"sub"`
	Role      string `json:"role"`
	SubjectID string `json:"subjectId,omitempty"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

type AuthService struct {
	db     *gorm.DB
	secret []byte
}

func NewAuthService(db *gorm.DB, secret string) *AuthService {
	if secret == "" {
		generated := make([]byte, 32)
		_, _ = rand.Read(generated)
		return &AuthService{db: db, secret: generated}
	}
	return &AuthService{db: db, secret: []byte(secret)}
}

func (s *AuthService) Login(identifier, password string) (*models.User, string, error) {
	var user models.User
	identifier = strings.TrimSpace(identifier)
	if err := s.db.Where("LOWER(subject_id) = LOWER(?)", identifier).First(&user).Error; err != nil {
		return nil, "", errors.New("invalid email or password")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, "", errors.New("invalid email or password")
	}
	subject := ""
	if user.SubjectID != nil {
		subject = *user.SubjectID
	}
	now := time.Now()
	token, err := s.sign(AuthClaims{UserID: user.ID, Role: user.Role, SubjectID: subject, IssuedAt: now.UnixNano(), ExpiresAt: now.Add(12 * time.Hour).Unix()})
	return &user, token, err
}

func (s *AuthService) CreateLinkedAccount(name, role, subjectID string) (string, string, error) {
	local := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		if r >= 'A' && r <= 'Z' {
			return r + ('a' - 'A')
		}
		return -1
	}, name)
	if local == "" {
		local = "user"
	}
	email := local + strings.ToLower(subjectID) + "@grgi.edu"
	prefix := "Student"
	if role == "professor" {
		prefix = "Teacher"
	}
	password := prefix + subjectID + "!"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	user := models.User{Name: name, Email: email, PasswordHash: string(hash), Role: role, SubjectID: &subjectID}
	if err := s.db.Create(&user).Error; err != nil {
		return "", "", err
	}
	return email, password, nil
}

func (s *AuthService) DeleteLinkedAccount(role, subjectID string) error {
	return s.db.Where("role = ? AND subject_id = ?", role, subjectID).Delete(&models.User{}).Error
}

func (s *AuthService) GetUser(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) UpdateUser(id uint, name, email string) (*models.User, error) {
	user, err := s.GetUser(id)
	if err != nil {
		return nil, err
	}
	trimmedName := strings.TrimSpace(name)
	if trimmedName != "" {
		user.Name = trimmedName
	}
	if strings.TrimSpace(email) != "" {
		user.Email = strings.ToLower(strings.TrimSpace(email))
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{"name": user.Name, "email": user.Email}
		if err := tx.Model(&models.User{}).Where("id = ?", user.ID).UpdateColumns(updates).Error; err != nil {
			return err
		}

		if trimmedName == "" || user.SubjectID == nil {
			return nil
		}
		switch user.Role {
		case "student":
			studentID, err := strconv.Atoi(*user.SubjectID)
			if err != nil {
				return errors.New("linked student id is invalid")
			}
			result := tx.Model(&models.Student{}).Where("id = ?", studentID).Update("name", trimmedName)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("linked student not found")
			}
		case "professor":
			result := tx.Model(&models.Professor{}).Where("id = ?", *user.SubjectID).Update("name", trimmedName)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("linked professor not found")
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) ChangePassword(id uint, current, next string) error {
	if len(next) < 8 {
		return errors.New("new password must be at least 8 characters")
	}
	user, err := s.GetUser(id)
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(current)) != nil {
		return errors.New("current password is incorrect")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(next), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.db.Model(user).Update("password_hash", string(hash)).Error
}

func (s *AuthService) DeleteUser(id uint, password string) error {
	user, err := s.GetUser(id)
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return errors.New("password is incorrect")
	}
	return s.db.Delete(user).Error
}

func (s *AuthService) Parse(token string) (*AuthClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token")
	}
	signed := parts[0] + "." + parts[1]
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("invalid token")
	}
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(signed))
	if !hmac.Equal(signature, mac.Sum(nil)) {
		return nil, errors.New("invalid token")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token")
	}
	var claims AuthClaims
	if json.Unmarshal(payload, &claims) != nil || claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("expired token")
	}
	var user models.User
	if err := s.db.Select("updated_at").First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("invalid session")
	}
	if claims.IssuedAt == 0 || claims.IssuedAt <= user.UpdatedAt.UnixNano() {
		return nil, errors.New("session has been logged out")
	}
	return &claims, nil
}

func (s *AuthService) RevokeUserTokens(userID uint) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("updated_at", time.Now()).Error
}

func (s *AuthService) sign(claims AuthClaims) (string, error) {
	header, _ := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	signed := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(signed))
	return signed + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func SubjectStudentID(claims *AuthClaims) (int, bool) {
	id, err := strconv.Atoi(claims.SubjectID)
	return id, err == nil
}
