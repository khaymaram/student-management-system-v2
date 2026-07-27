package models

import "time"

type Course struct {
	ID uint `gorm:"column:id;primaryKey" json:"courseId"`

	Title string `gorm:"column:title" json:"title"`

	Code string `gorm:"column:code" json:"code"`

	Credits int `gorm:"column:credits" json:"credits"`

	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`

	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

