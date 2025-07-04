package seeders

import (
	"encoding/json"
	"go-gin-hexagonal/internal/adapter/database/model"
	"io"
	"log"
	"os"

	"gorm.io/gorm"
)

func UserSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("pkg/database/json/user.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var users []model.User
	if err := json.Unmarshal(jsonData, &users); err != nil {
		return err
	}

	for _, user := range users {
		isData := db.Find(&user, "email = ? OR username = ?", user.Email, user.Username).RowsAffected
		if isData == 0 {
			if err := db.Create(&user).Error; err != nil {
				log.Printf("error seeding user: %v", err)
				return err
			}
		}
	}

	return nil
}
