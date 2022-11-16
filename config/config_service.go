package config

import (
	"context"
	"errors"
	"github.com/deathsgun/art/auth"
	"github.com/deathsgun/art/config/model"
	"github.com/deathsgun/art/crypt"
	"github.com/deathsgun/art/di"
	"gorm.io/gorm"
)

type IConfigService interface {
	GetConfigsForUser(ctx context.Context) ([]model.ProviderConfig, error)
	GetConfiguredProviderNames(ctx context.Context) ([]string, error)
	SaveProviderConfig(ctx context.Context, config *model.ProviderConfig) error
	GetConfig(ctx context.Context, provider string) (*model.ProviderConfig, error)
	DeleteConfig(ctx context.Context, provider string) error
}

type service struct {
}

func (s *service) DeleteConfig(ctx context.Context, provider string) error {
	db := di.Instance[*gorm.DB]("database")
	result := db.Where(&model.ProviderConfig{User: auth.Session(ctx).Id(), Provider: provider}).Delete(&model.ProviderConfig{})
	return result.Error
}

func (s *service) GetConfiguredProviderNames(ctx context.Context) ([]string, error) {
	db := di.Instance[*gorm.DB]("database")
	var configs []model.ProviderConfig
	result := db.Select("provider", "user").Where("user = ?", auth.Session(ctx).Id()).Find(&configs)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	var names []string
	for _, config := range configs {
		names = append(names, config.Provider)
	}
	return names, nil
}

func (s *service) GetConfigsForUser(ctx context.Context) ([]model.ProviderConfig, error) {
	db := di.Instance[*gorm.DB]("database")
	var configs []model.ProviderConfig
	result := db.Where(&model.ProviderConfig{User: auth.Session(ctx).Id()}).Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}
	return configs, nil
}

func (s *service) SaveProviderConfig(ctx context.Context, config *model.ProviderConfig) error {
	db := di.Instance[*gorm.DB]("database")

	pw, err := di.Instance[crypt.ICryptService]("crypt").EncryptString(config.Password)
	if err != nil {
		return err
	}
	config.Password = pw
	config.User = auth.Session(ctx).Id()
	result := db.First(&model.ProviderConfig{}, &model.ProviderConfig{User: auth.Session(ctx).Id(), Provider: config.Provider})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		result = db.Create(config)
		return result.Error
	}
	result = db.Where(&model.ProviderConfig{User: auth.Session(ctx).Id(), Provider: config.Provider}).Updates(config)
	return result.Error
}

func (s *service) GetConfig(ctx context.Context, provider string) (*model.ProviderConfig, error) {
	db := di.Instance[*gorm.DB]("database")
	config := &model.ProviderConfig{}
	result := db.Where(&model.ProviderConfig{User: auth.Session(ctx).Id(), Provider: provider}).First(config)
	if result.Error != nil {
		return nil, result.Error
	}
	pw, err := di.Instance[crypt.ICryptService]("crypt").DecryptString(config.Password)
	if err != nil {
		return nil, err
	}
	config.Password = pw
	return config, nil
}

func New() IConfigService {
	return &service{}
}
