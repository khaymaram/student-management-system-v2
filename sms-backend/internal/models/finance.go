package models

import "time"

type Finance struct {
	StudentID int `gorm:"column:student_id;primaryKey;autoIncrement:false;not null"`

	Student *Student `gorm:"foreignKey:StudentID;references:ID;constraint:OnDelete:CASCADE" json:"student,omitempty"`

	Tuition float64 `gorm:"column:tuition;type:decimal(10,2);not null"`

	Scholarship float64 `gorm:"column:scholarship;type:decimal(10,2);default:0"`

	Paid float64 `gorm:"column:paid;type:decimal(10,2);default:0"`

	IsInState bool `gorm:"column:is_in_state;not null"`

	CreatedAt time.Time

	UpdatedAt time.Time
}

func (Finance) TableName() string {
	return "Finances"
}
