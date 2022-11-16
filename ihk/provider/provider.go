package provider

import (
	"context"
	"errors"
	"github.com/deathsgun/art/config"
	"github.com/deathsgun/art/config/model"
	"github.com/deathsgun/art/di"
	"github.com/deathsgun/art/ihk"
	"github.com/deathsgun/art/provider"
	"github.com/deathsgun/art/report"
	"github.com/deathsgun/art/utils"
	"time"
)

type impl struct {
}

func (i *impl) Id() string {
	return "PROVIDER_IHK"
}

func (i *impl) Logo() string {
	return "ihk-logo.png"
}

func (i *impl) Capabilities() []provider.Capability {
	return []provider.Capability{
		provider.TypeExport,
		provider.Configurable,
		provider.ConfigServer,
		provider.ConfigUsername,
		provider.ConfigPassword,
		provider.ConfigDepartment,
		provider.ConfigInstructorEmail,
		provider.ConfigSendDirectly,
	}
}

func (i *impl) ValidateConfig(ctx context.Context, config *model.ProviderConfig) error {
	ihkService := di.Instance[ihk.IIHKService]("ihkService")
	token, err := ihkService.Login(ctx, config.Server, config.Username, config.Password)
	if err != nil {
		return err
	}
	return ihkService.Logout(ctx, config.Server, token)
}

func (i *impl) Export(ctx context.Context, rep *report.Report) ([]byte, error) {
	ihkService := di.Instance[ihk.IIHKService]("ihkService")
	configService := di.Instance[config.IConfigService]("configService")
	conf, err := configService.GetConfig(ctx, i.Id())
	if err != nil {
		return nil, err
	}
	if conf.InstructorEmail == "" {
		return nil, errors.New("INSTRUCTOR_EMAIL_EMPTY")
	}
	if conf.Department == "" {
		return nil, errors.New("DEPARTMENT_EMPTY")
	}

	tokens, err := ihkService.Login(ctx, conf.Server, conf.Username, conf.Password)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = ihkService.Logout(ctx, conf.Server, tokens)
	}()

	ihkReport, err := ihkService.CreateNewReport(ctx, conf.Server, conf.InstructorEmail, conf.Department, tokens)
	if err != nil {
		return nil, err
	}

	entries, err := rep.Format(ctx, true)
	if err != nil {
		return nil, err
	}

	for category, bytes := range entries {
		switch category {
		case report.Activity:
			ihkReport.ContentActivity = string(bytes)
		case report.Subjects:
			ihkReport.ContentSubjects = string(bytes)
		case report.Training:
			ihkReport.ContentTraining = string(bytes)
		}
	}

	return nil, ihkService.SaveReport(ctx, conf.Server, ihkReport, tokens)
}

func (i *impl) GetStartDate(ctx context.Context) (time.Time, error) {
	configService := di.Instance[config.IConfigService]("configService")
	conf, err := configService.GetConfig(ctx, i.Id())
	if err != nil {
		return time.Now(), err
	}
	ihkService := di.Instance[ihk.IIHKService]("ihkService")
	tokens, err := ihkService.Login(ctx, conf.Server, conf.Username, conf.Password)
	if err != nil {
		return time.Now(), err
	}
	defer func() {
		err := ihkService.Logout(ctx, conf.Server, tokens)
		if err != nil {
			return
		}
	}()
	ihkReport, err := ihkService.CreateNewReport(ctx, conf.Server, conf.InstructorEmail, conf.Department, tokens)
	if err != nil {
		return time.Now(), err
	}

	err = ihkService.CancelReport(ctx, conf.Server, ihkReport, tokens)
	if err != nil {
		return time.Now(), err
	}

	return utils.LeapToNearestMonday(ihkReport.StartDate), nil
}

func (i *impl) ContentType() string {
	return "none"
}

func NewProvider() provider.ExportProvider {
	return &impl{}
}
