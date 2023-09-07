package database

import (
	"log"
	"strconv"
	"strings"

	"example.com/m/models"
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

// connectDb
func ConnectDb() {

	db, err := gorm.Open(sqlite.Open("locations.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Silent)
	log.Println("running migrations")
	db.AutoMigrate(&models.Location{})
	//	DB.AutoMigrate(&model.Product{}, &model.User{})
	DB = Dbinstance{
		Db: db,
	}
}
func Paginate(r *fiber.Ctx) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		page, _ := strconv.Atoi(r.Query("page"))
		if page <= 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(r.Query("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
func LocationByName(nameFilter string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("Logradouro LIKE ?", "%"+strings.ToUpper(nameFilter)+"%")
	}
}
