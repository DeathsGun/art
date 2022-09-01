package provider

import (
	"fmt"
	"os"
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
	outputFile := outputDir + "\\" + startDate + ".txt"
	content := ""
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
		content += "- " + v.Text + " \n"
	}

	textFile.WriteString(content) //NOCH NICHT FERTIG!!!! Keine Fiesen Kommentare erst Mittwoch!

	if err != nil {
		return err
	}

	return err
}

func NewTextProvider() ExportProvider {
	return &textProvider{}
}
