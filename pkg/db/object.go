package db

type Project struct {
	ID          string `gorm:"column:id;type:varchar(255);primary_key" json:"id"`
	ProjectInfo `json:"detail"`
}

type ProjectInfo struct {
	Project string `gorm:"column:project;type:varchar(255)" json:"project"`
	Product string `gorm:"column:product;type:varchar(255)" json:"product"`
	Profile string `gorm:"column:profile;type:varchar(255)" json:"profile"`
	Version string `gorm:"column:version;type:varchar(255)" json:"version"`
}
