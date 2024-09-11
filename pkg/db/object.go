package db

type Project struct {
	ID      string `gorm:"column:id;type:varchar(255);primary_key" json:"id"`
	Project string `gorm:"column:project;type:varchar(255)" json:"project"`
	Name    string `gorm:"column:name;type:varchar(255)" json:"name"`
	Profile string `gorm:"column:profile;type:varchar(255)" json:"profile"`
}
