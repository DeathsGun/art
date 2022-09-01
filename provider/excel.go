package provider

import (
	"fmt"

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

func (e *excelProvider) Export(report *Report, startDate string, templateFile string, outputDir string) error {
	outputFile := outputDir + "\\" + startDate + ".xlsx"
	content := ""
	excelFile, err := excelize.OpenFile(templateFile)
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

	for _, v := range report.Entries {
		content += "- " + v.Text + " \n"
	}

	excelFile.SetCellValue(excelFile.GetSheetName(0), "A42", content)
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
