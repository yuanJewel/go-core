package db

import "time"

type User struct {
	ID            int       `gorm:"column:id;type:int;primary_key;AUTO_INCREMENT" json:"id"`
	Name          string    `gorm:"column:name;type:varchar(255);uniqueIndex" json:"name"`
	Passwd        string    `gorm:"column:passwd;type:varchar(255)" json:"passwd"`
	GoogleSecret  string    `gorm:"column:google_secret;type:varchar(255)" json:"google_secret"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
	LastLoginTime time.Time `gorm:"column:last_login_time;type:datetime" json:"last_login_time"`
}

type AssetRecord struct {
	ID           int       `gorm:"column:id;type:int;primary_key;AUTO_INCREMENT" json:"id"`
	TranceId     string    `gorm:"column:trance_id;type:varchar(255)" json:"trance_id"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
	UpdateUserId int       `gorm:"column:update_user_id;type:int" json:"update_user_id"`
	Url          string    `gorm:"column:url;type:varchar(255)" json:"url"`
	Method       string    `gorm:"column:method;type:varchar(255)" json:"method"`
	Body         string    `gorm:"column:body;type:varchar(255)" json:"body"`
	Context      string    `gorm:"column:context;type:text" json:"context"`
}

type TableAffect struct {
	ID           int       `gorm:"column:id;type:int;primary_key;AUTO_INCREMENT" json:"id"`
	TranceId     string    `gorm:"column:trance_id;type:varchar(255)" json:"trance_id"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
	UpdateUserId int       `gorm:"column:update_user_id;type:int" json:"update_user_id"`
	Table        string    `gorm:"column:table;type:varchar(255)" json:"table"`
	PrimaryId    string    `gorm:"column:primary_id;type:varchar(255)" json:"primary_id"`
	Action       string    `gorm:"column:action;type:varchar(255)" json:"action"`
}
