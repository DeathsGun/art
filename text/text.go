package text

import (
	"errors"
	"github.com/deathsgun/art/provider"
	"os"
	"path/filepath"
	"time"
)

type textProvider struct {
}

func (t *textProvider) Name() string {
	return "text"
}

func (t *textProvider) ValidateLogin(username string, password string) (string, string, error) {
	return username, password, nil
}

func (t *textProvider) NeedsLogin() bool {
	return false
}

func (t *textProvider) Export(report *provider.Report, startDate time.Time, outputDir string, printDates bool) error {
	if outputDir == "" {
		return errors.New("output dir required")
	}
	_ = os.MkdirAll(outputDir, os.ModePerm)

	file, err := os.OpenFile(filepath.Join(outputDir, startDate.Format("2006-01-02")+".txt"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	buffers, err := report.Format(printDates)
	if err != nil {
		return err
	}
	for category, bytes := range buffers {
		_, err = file.WriteString(category.Text() + ":\n")
		if err != nil {
			return err
		}
		if len(bytes) == 0 {
			continue
		}
		_, err = file.Write(bytes)
		if err != nil {
			return err
		}
		_, err = file.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func NewTextProvider() provider.ExportProvider {
	return &textProvider{}
}
