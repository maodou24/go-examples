package modle

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name string `gorm:"column:name; NOT NULL; unique"`
	Age  int    `gorm:"column:age"`
}
