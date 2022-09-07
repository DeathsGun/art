package provider

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type textProvider struct {
}

func (t *textProvider) Name() string {
	return "text"
}

func (t *textProvider) ValidateLogin(username string, password string) error {
	return nil
}

func (t *textProvider) NeedsLogin() bool {
	return false
}

func (t *textProvider) Export(report *Report, startDate time.Time, outputDir string) error {
	if outputDir == "" {
		return errors.New("output dir required")
	}
	_ = os.MkdirAll(outputDir, os.ModePerm)

	file, err := os.OpenFile(filepath.Join(outputDir, startDate.Format("2006-01-02")+".txt"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	var categories []Category
	var dates []string
	for _, v := range report.Entries {
		if !Contains(categories, v.Category) {
			_, err = file.WriteString(v.Category.Text() + "\n")
			if err != nil {
				return err
			}
			categories = append(categories, v.Category)
		}
		date := v.Date.Format("Monday 02.01.2006:")
		if !Contains(dates, date) {
			_, err = file.WriteString("\t" + date + "\n")
			if err != nil {
				return err
			}
			dates = append(dates, date)
		}

		_, err = file.WriteString(fmt.Sprintf("\t\t- %s\n", v.Text))
		if err != nil {
			return err
		}
	}
	return nil
}

func Contains[T comparable](array []T, val T) bool {
	for _, t := range array {
		if t == val {
			return true
		}
	}
	return false
}

func NewTextProvider() ExportProvider {
	return &textProvider{}
}
