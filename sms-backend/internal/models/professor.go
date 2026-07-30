package models

import "time"

type Professor struct {
	ID string `gorm:"column:id;type:nvarchar(50);primaryKey" json:"id"`

	Name string `gorm:"column:name;not null" json:"name"`

	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`

	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`

	Courses []Course `gorm:"foreignKey:ProfessorID;references:ID" json:"courses,omitempty"`
}
func (Professor) TableName() string {
	return "dbo.Professors"
}