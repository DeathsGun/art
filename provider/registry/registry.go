package registry

import (
	"github.com/deathsgun/art/provider"
)

// ImportProviders is the global store where provider.ImportProvider get registered
var ImportProviders = []provider.ImportProvider{}

// ExportProviders is the global store where provider.ExportProvider get registered
var ExportProviders = []provider.ExportProvider{}
