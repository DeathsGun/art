package provider

import (
	"context"
	"github.com/deathsgun/art/report"
	"time"
)

type Capability int

const (
	Configurable Capability = iota
	ConfigServer
	ConfigUsername
	ConfigPassword
	TypeImport
	TypeExport
)

// Provider has just methods required for both ImportProvider and ExportProvider
type Provider interface {
	// Id is the translation key of the provider
	Id() string
	Logo() string
	Capabilities() []Capability
}

type ImportProvider interface {
	Provider
	// Import gives the provider a monday based date which has been set 00:00h,
	// so they can use this date to import data for the whole week instead of some random
	// days
	Import(ctx context.Context, monday time.Time) ([]report.Entry, error)
}

type ExportProvider interface {
	Provider
	// Export receives the Report from export.HandleExport with entries filled from all the
	// configured ImportProvider
	Export(ctx context.Context, report *report.Report) ([]byte, error)
}
