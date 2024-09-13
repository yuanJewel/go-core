package cmdb

import (
	"gorm.io/gorm"
)

type Service interface {
	HealthCheck() error
	Setup([]interface{}) error

	GetTables() ([]string, error)
	HasTable(string) bool

	// General operations for manipulating databases
	AddItem(interface{}, int64) (*gorm.DB, error)
	UpdateItem(interface{}, interface{}, int64) (*gorm.DB, error)
	DeleteItem(interface{}, int64) (*gorm.DB, error)
	GetItems(interface{}, interface{}) (bool, error)
	GetItemsOrder(interface{}, interface{}, string) (bool, error)
	GetItemsFromSlice(interface{}, interface{}, ...interface{}) (bool, error)
	GetItemsFromSliceOrder(interface{}, string, interface{}, ...interface{}) (bool, error)
	GetItemsFromDataAndSlice(interface{}, string, interface{}, ...interface{}) (bool, error)
	GetItemsFromDataAndSliceOrder(interface{}, string, string, interface{}, ...interface{}) (bool, error)
	GetItem(interface{}, interface{}) (bool, error)

	// Customized operation
	GetItemsByIds(interface{}, interface{}, string) (bool, error)
	GetItemsByIdsOrder(interface{}, interface{}, string, string) (bool, error)
	GetItemsByIdsAndSlices(interface{}, interface{}, string, string, string) (bool, error)
}
