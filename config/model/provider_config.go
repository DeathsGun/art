package model

import "gorm.io/gorm"

type ProviderConfig struct {
	gorm.Model
	Provider string
	User     string
	Server   string
	Username string
	Password string
	Token    string
}
