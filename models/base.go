package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/satori/go.uuid"
	"log"
	"os"
	"time"
)

var db *gorm.DB // database
var username, password, dbName, dbHost, dbPort string

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}

	username = os.Getenv("db_user")
	password = os.Getenv("db_pass")
	dbName = os.Getenv("db_name")
	dbHost = os.Getenv("db_host")
	dbPort = os.Getenv("db_port")

	migrateDatabase()
}

// Datebase migration
func migrateDatabase() {
	db := GetDB()

	db.Debug().AutoMigrate(
		&Attendee{},
		&User{},
	)

	// Migration scripts
	//db.Model(&Attendee{}).AddForeignKey("parent_id", "attendees(id)", "SET NULL", "RESTRICT")
}

func GetDB() *gorm.DB {
	dbUri := fmt.Sprintf("postgres://%v@%v:%v/%v?sslmode=disable&password=%v", username, dbHost, dbPort, dbName, password)

	// Making connection to the database
	db, err := gorm.Open("postgres", dbUri)
	if err != nil {
		log.Println(err)
	}

	return db
}
