package models

// Major is a selectable academic program. Rows are managed as lookup data.
type Major struct {
	ID   int    `gorm:"column:id;primaryKey" json:"id"`
	Name string `gorm:"column:name;type:nvarchar(100);uniqueIndex;not null" json:"name"`
}

func (Major) TableName() string {
	return "dbo.Majors"
}
