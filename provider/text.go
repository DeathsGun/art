package provider

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
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

func (t *textProvider) Export(report *Report, startDate string, templateFile string, outputDir string) error {
	fileName := startDate + ".txt"
	outputFile := filepath.Join(outputDir, fileName)
	buffer := &bytes.Buffer{}

	textFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := textFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for _, v := range report.Entries {
		buffer.WriteString(fmt.Sprintf("- %s\n", v.Text))
	}

	textFile.WriteString(buffer.String())

	if err != nil {
		return err
	}

	return err
}

func NewTextProvider() ExportProvider {
	return &textProvider{}
}
