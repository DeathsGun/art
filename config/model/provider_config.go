package model

import "gorm.io/gorm"

type ProviderConfig struct {
	gorm.Model
	Provider        string
	User            string
	Server          string
	Username        string
	Password        string
	InstructorEmail string
	Department      string
	SendDirectly    bool
	Token           string
}
