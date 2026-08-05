// Package models defines the database entity schemas for the backend.
// Student maps the database row to JSON fields used by the API.
package models

import "time"

// Student is the database model used by GORM and JSON responses.
type Student struct {
	ID int `gorm:"column:id;primaryKey;autoIncrement:false" json:"studentId"`

	Name string `gorm:"column:name" json:"name"`

	Grade int64 `gorm:"column:grade" json:"grade"`

	// GPA is a derived, credit-weighted value calculated from graded enrollments.
	GPA *float64 `gorm:"column:gpa" json:"gpa"`

	Finance *Finance `gorm:"foreignKey:StudentID;references:ID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`

	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (Student) TableName() string {
	return "dbo.Students"
}