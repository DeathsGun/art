package provider

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/user"
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

func (e *excelProvider) Export(report *Report, startDate time.Time, outputDir string) error {
	if outputDir == "" {
		return errors.New("output dir required")
	}
	_ = os.MkdirAll(outputDir, os.ModePerm)
	outputFile := filepath.Join(outputDir, startDate.Format("2006-01-02")+".xlsx")
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

	mondayDiff := time.Monday - startDate.Weekday()
	monday := startDate.Add(time.Hour * time.Duration(24*mondayDiff))

	sheet := excelFile.GetSheetName(0)
	// Begin
	err = excelFile.SetCellValue(sheet, "E7", monday.Format("02.01.2006"))
	if err != nil {
		return err
	}
	err = excelFile.SetCellValue(sheet, "G7", monday.Add(4*24*time.Hour).Format("02.01.2006"))
	if err != nil {
		return err
	}
	// Content
	err = excelFile.SetCellValue(sheet, "A42", buffer.Bytes())
	if err != nil {
		return err
	}

	// User
	u, err := user.Current()
	if err != nil {
		return err
	}
	err = excelFile.SetCellValue(sheet, "G2", u.Name)
	if err != nil {
		return err
	}

	if err := excelFile.SaveAs(outputFile); err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func NewExcelProvider() ExportProvider {
	return &excelProvider{}
}
