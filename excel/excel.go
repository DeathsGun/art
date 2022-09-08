package excel

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"github.com/deathsgun/art/provider"
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

func (e *excelProvider) ValidateLogin(username string, password string) (string, string, error) {
	return username, password, nil
}

func (e *excelProvider) NeedsLogin() bool {
	return false
}

//go:embed Template.xlsx
var template []byte

// Export is
func (e *excelProvider) Export(report *provider.Report, startDate time.Time, outputDir string, printDate bool) error {
	if outputDir == "" {
		return errors.New("output dir required")
	}

	buffers, err := report.Format(printDate)
	if err != nil {
		return err
	}

	_ = os.MkdirAll(outputDir, os.ModePerm)
	outputFile := filepath.Join(outputDir, startDate.Format("2006-01-02")+".xlsx")

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

	//Content
	err = excelFile.SetCellValue(sheet, "A10", buffers[provider.ACTIVITY])
	if err != nil {
		return err
	}

	err = excelFile.SetCellValue(sheet, "A30", buffers[provider.TRAINING])
	if err != nil {
		return err
	}

	err = excelFile.SetCellValue(sheet, "A42", buffers[provider.SUBJECTS])
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

func NewExcelProvider() provider.ExportProvider {
	return &excelProvider{}
}
