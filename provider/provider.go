package provider

import "time"

type Provider interface {
	Name() string
	ValidateLogin(username string, password string) error
	NeedsLogin() bool
}

type ImportProvider interface {
	Provider
	Import(startDate string) ([]*Entry, error)
}

type ExportProvider interface {
	Provider
	Export(report *Report, startDate time.Time, outputDir string) error
}
