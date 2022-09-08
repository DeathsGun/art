package provider

import (
	"errors"
	"time"
)

var (
	ErrNoLoginConfigured = errors.New("no login configured")
)

type Provider interface {
	Name() string
	ValidateLogin(username string, password string) (string, string, error)
	NeedsLogin() bool
}

type ImportProvider interface {
	Provider
	Import(startDate time.Time) ([]*Entry, error)
}

type ExportProvider interface {
	Provider
	Export(report *Report, startDate time.Time, outputDir string, printDates bool) error
}
