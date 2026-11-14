package main

import (
	"api/config"
	"api/db"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := gorm.Open(mysql.Open(config.GetDSN(cfg)), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("✓ Database connected successfully")
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("Failed to prepare database schema: %v", err)
	}
	if err := db.SeedInitialData(database); err != nil {
		log.Fatalf("Failed to seed starter data: %v", err)
	}
	if err := db.SeedAdmin(database, cfg.AdminName, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		log.Fatalf("Failed to seed admin: %v", err)
	}
	log.Println("✓ Seeding completed successfully")
}
