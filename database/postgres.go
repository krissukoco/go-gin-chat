package database

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres() (*gorm.DB, error) {
	dbHost, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		return nil, errors.New("POSTGRES_HOST is not set")
	}
	dbUser, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		return nil, errors.New("POSTGRES_USER is not set")
	}
	dbPassword, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		return nil, errors.New("POSTGRES_PASSWORD is not set")
	}
	dbName, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		return nil, errors.New("POSTGRES_DB is not set")
	}
	dbPort, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		return nil, errors.New("POSTGRES_PORT is not set")
	}
	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbUser, dbPassword, dbName, dbPortInt)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}
