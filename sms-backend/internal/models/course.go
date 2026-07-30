package models

import "time"

type Course struct {
	// ID uint `gorm:"column:id;primaryKey" json:"courseId"`
	Code string `gorm:"column:code;primaryKey;type:nvarchar(256)" json:"code"`

	Title string `gorm:"column:title" json:"title"`

	Credits int `gorm:"column:credits" json:"credits"`

	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`

	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`

	ProfessorID string     `gorm:"column:professor_id;type:nvarchar(50)" json:"professorId"`
	Professor   *Professor `gorm:"foreignKey:ProfessorID;references:ID" json:"professor,omitempty"`
}

func (Course) TableName() string {
	return "dbo.courses"
}
