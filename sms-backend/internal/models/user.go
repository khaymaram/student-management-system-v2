package models

import "time"

type User struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"column:name;not null" json:"name"`
	Email        string    `gorm:"column:email;type:nvarchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;type:nvarchar(255);not null" json:"-"`
	Role         string    `gorm:"column:role;type:nvarchar(20);not null" json:"role"`
	SubjectID    *string   `gorm:"column:subject_id;type:nvarchar(50);uniqueIndex" json:"subjectId"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (User) TableName() string { return "dbo.Users" }
