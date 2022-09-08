package provider

import (
	"errors"
	"time"
)

var (
	ErrNoLoginConfigured = errors.New("no login configured")
)

// Provider has just methods required for both ImportProvider and ExportProvider
type Provider interface {
	// Name is the lowercase name of the provider
	Name() string
	// ValidateLogin allows the provider to check the credentials and return
	// tokens if required instead of username's and passwords
	ValidateLogin(username string, password string) (string, string, error)
	// NeedsLogin allows the provider specify whether they need a login or not
	NeedsLogin() bool
}

type ImportProvider interface {
	Provider
	// Import gives the provider a monday based date which has been set 00:00h,
	// so they can use this date to import data for the whole week instead of some random
	// days
	Import(startDate time.Time) ([]*Entry, error)
}

type ExportProvider interface {
	Provider
	// Export receives the Report from export.HandleExport with entries filled from all the
	// configured ImportProvider
	Export(report *Report, startDate time.Time, outputDir string, printDates bool) error
}
