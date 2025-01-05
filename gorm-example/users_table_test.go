package gorm_example

import (
	"fmt"
	"testing"

	"github.com/maodou24/gorm-example/modle"
	"github.com/maodou24/gorm-example/pg"
)

func TestCreateTable(t *testing.T) {
	if err := pg.GetPgTestDb().AutoMigrate(&modle.User{}); err != nil {
		t.Error(err)
	}
}

func TestUserTableInsert(t *testing.T) {
	user := modle.User{Name: "maodou", Age: 18}
	pg.GetPgTestDb().Create(&user)
}

func TestUserTableQueryFirst(t *testing.T) {
	// query
	var result modle.User
	pg.GetPgTestDb().First(&result)

	fmt.Printf("%+v", result)
}

func TestUserTableQuery(t *testing.T) {
	// query
	user := modle.User{Name: "maodou"}
	pg.GetPgTestDb().Find(&user)

	fmt.Printf("%+v", user)
}

func TestUserTableDelete(t *testing.T) {
	deleteUser := modle.User{Name: "maodou"}
	if err := pg.GetPgTestDb().Where("name = ?", deleteUser.Name).Delete(&modle.User{}).Error; err != nil {
		t.Error(err)
	}
}

func TestUserTableDeleteByAge(t *testing.T) {
	pg.GetPgTestDb().Where("age = ?", 18).Delete(&modle.User{})
}
