package dto

import "github.com/deathsgun/art/config/model"

type ProviderConfig struct {
	Server          string `json:"server"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Department      string `json:"department"`
	InstructorEmail string `json:"instructor"`
}

// ToDto TODO don't send password LOL
func ToDto(m *model.ProviderConfig) *ProviderConfig {
	return &ProviderConfig{
		Server:          m.Server,
		Username:        m.Username,
		Password:        "dummypassword",
		Department:      m.Department,
		InstructorEmail: m.InstructorEmail,
	}
}

func ToModel(d *ProviderConfig) *model.ProviderConfig {
	if d.Password == "dummypassword" {
		d.Password = ""
	}
	return &model.ProviderConfig{
		Server:          d.Server,
		Username:        d.Username,
		Password:        d.Password,
		Department:      d.Department,
		InstructorEmail: d.InstructorEmail,
	}
}
