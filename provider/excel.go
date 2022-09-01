package provider

type excelProvider struct {
}

func (e *excelProvider) Name() string {
	return ""
}

func (e *excelProvider) ValidateLogin(username string, password string) error {
	return nil
}

func (e *excelProvider) NeedsLogin() bool {
	return false
}

func (e *excelProvider) Export(report *Report, startDate string, templateFile string, outputDir string) error {
	return nil
}

func NewExcelProvider() ExportProvider {
	return &excelProvider{}
}
