package provider

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

type excelProvider struct {
}

func (e *excelProvider) Name() string {
	return "excel"
}

func (e *excelProvider) ValidateLogin(username string, password string) error {
	return nil
}

func (e *excelProvider) NeedsLogin() bool {
	return false
}

//go:embed Template.xlsx
var template []byte

func (e *excelProvider) Export(report *Report, startDate string, cell string, outputDir string) error {
	if outputDir == "" {
		return errors.New("output dir required")
	}
	if startDate == "" {
		startDate = time.Now().Format("2006-01-02")
	}
	outputFile := filepath.Join(outputDir, startDate+".xlsx")
	buffer := &bytes.Buffer{}

	excelFile, err := excelize.OpenReader(bytes.NewReader(template))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := excelFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var dates []string
	for _, v := range report.Entries {
		date := v.Date.Format("Monday 02.01.2006:")
		if !Contains(dates, date) {
			_, err = buffer.WriteString("\t" + date + "\n")
			if err != nil {
				return err
			}
			dates = append(dates, date)
		}

		_, err = buffer.WriteString(fmt.Sprintf("\t\t- %s\n", v.Text))
		if err != nil {
			return err
		}
	}

	if cell == "" {
		cell = "A42"
	}

	err = excelFile.SetCellValue(excelFile.GetSheetName(0), cell, buffer.Bytes())
	if err != nil {
		return err
	}
	fmt.Println("Excel Export succeed!")
	if err := excelFile.SaveAs(outputFile); err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func NewExcelProvider() ExportProvider {
	return &excelProvider{}
}
