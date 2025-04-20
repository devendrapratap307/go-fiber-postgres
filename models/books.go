package models

import "gorm.io/gorm"

type Books struct {
	ID        uint    `json:"id" gorm:"primaryKey; autoIncrement"`
	Title     *string `json:"title" gorm:"not null"`
	Author    *string `json:"author" gorm:"not null"`
	Year      *int    `json:"year" gorm:"not null"`
	Publisher *string `json:"publisher" gorm:"not null"`
}

func MigrateBooks(db *gorm.DB) error {
	if err := db.AutoMigrate(&Books{}); err != nil {
		return err
	}
	return nil
}
