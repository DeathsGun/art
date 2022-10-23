package dto

import "github.com/deathsgun/art/config/model"

type ProviderConfig struct {
	Server   string `json:"server"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// ToDto TODO don't send password LOL
func ToDto(m *model.ProviderConfig) *ProviderConfig {
	return &ProviderConfig{
		Server:   m.Server,
		Username: m.Username,
		Password: m.Password,
	}
}

func ToModel(d *ProviderConfig) *model.ProviderConfig {
	return &model.ProviderConfig{
		Server:   d.Server,
		Username: d.Username,
		Password: d.Password,
	}
}
