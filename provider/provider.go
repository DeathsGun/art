package provider

var ImportProviders = map[string]ImportProvider{}
var ExportProviders = map[string]ExportProvider{"test": NewTestProvider()}

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
