package provider

var ImportProviders []ImportProvider
var ExportProviders = []ExportProvider{NewTestProvider()}

type Provider interface {
	Name() string
	ValidateLogin(username string, password string) error
	NeedsLogin() bool
}

type ImportProvider interface {
	Provider
	Import() ([]*Entry, error)
}

type ExportProvider interface {
	Provider
	Export(report *Report) error
}
