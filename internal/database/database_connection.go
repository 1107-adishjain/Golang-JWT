package database

// in this file we will write the logic of connecting to DB.
import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBinitialize opens a GORM DB using the provided URL/DSN and returns the *gorm.DB handle.
func DBinitialize(url string) (*gorm.DB, error) {
	gdb, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return gdb, nil
}

// DBClose closes the underlying *sql.DB for the given *gorm.DB.
func DBClose(gdb *gorm.DB) error {
	sqlDB, err := gdb.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
