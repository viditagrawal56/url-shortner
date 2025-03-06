package db

import (
	"fmt"
	"log"

	"github.com/viditagrawal56/url-shortner/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	db *gorm.DB
}

func New(cfg config.DatabaseConfig) (*Database, error) {
	gormdb, err := gorm.Open(postgres.Open(cfg.ConnectionStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db, err := gormdb.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB instance: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	//Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{db: gormdb}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB for closing: %w", err)
	}

	return sqlDB.Close()
}

func (d *Database) AutoMigrate(models ...interface{}) error {
	return d.db.AutoMigrate(models...)
}

func (d *Database) ResetAndMigrate(models ...interface{}) error {
	// Drop tables in reverse order to handle foreign key constraints
	for i := len(models) - 1; i >= 0; i-- {
		if err := d.db.Migrator().DropTable(models[i]); err != nil {
			return fmt.Errorf("failed to drop table for model %T: %w", models[i], err)
		}
		log.Printf("Dropped table for model %T", models[i])
	}

	// Run migrations for all models
	if err := d.db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("Successfully re-created and migrated all tables")

	return nil
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}
