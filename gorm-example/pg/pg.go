package pg

import (
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDb *gorm.DB

func init() {
	dsn := "host=localhost user=postgres password=maodou24 dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}

	testDb = db
}

func GetPgTestDb() *gorm.DB {
	return testDb
}
