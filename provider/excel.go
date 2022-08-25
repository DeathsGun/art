package provider

type excelProvider struct {
}

func (e *excelProvider) Name() string {
	return ""
}

func (t *excelProvider) ValidateLogin(username string, password string) error {
	return nil
}

func (e *excelProvider) NeedsLogin() bool {
	return false
}

func (t *excelProvider) Export(report *Report) error {
	return nil
}

func NewExcelProvider() ExportProvider {
	return &excelProvider{}
}
