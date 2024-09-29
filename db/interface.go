package db

import (
	"gorm.io/gorm"
)

type Service interface {
	HealthCheck() error
	Setup([]interface{}) error

	GetTables() ([]string, error)
	HasTable(string) bool

	// Query operations
	Preload(string, ...interface{}) Service
	Joins(string, ...interface{}) Service
	Where(string, ...interface{}) Service
	Search(string, string) Service
	Order(string) Service
	Limit(int) Service
	OffsetPages(page int) Service

	// General operations for manipulating databases
	AddItem(interface{}, int64) (*gorm.DB, error)
	UpdateItem(interface{}, interface{}, int64) (*gorm.DB, error)
	DeleteItem(interface{}, int64) (*gorm.DB, error)
	GetItems(interface{}, interface{}) (int64, error)
	GetItem(interface{}, interface{}) (bool, error)
}
