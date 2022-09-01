package export

import (
	"fmt"
	"github.com/deathsgun/art/provider"
	"os"
)

func HandleExport(prov string, start string, temp string, output string) {
	if prov == "" {
		println("a provider for export is required")
		os.Exit(1)
	}

	var exportProvider provider.ExportProvider = nil
	for _, p := range provider.ExportProviders {
		if p.Name() == prov {
			exportProvider = p
		}
	}

	if exportProvider == nil {
		println("provider %s not found", prov)
		os.Exit(1)
	}

	report := &provider.Report{
		Id:      0,
		Entries: []*provider.Entry{},
	}

	for _, iprov := range provider.ImportProviders {
		entries, err := iprov.Import()
		if err != nil {
			fmt.Printf("Skipping import provider %s because it errored: %v", iprov.Name(), entries)
			continue
		}
		report.Entries = append(report.Entries, entries...)
	}

	err := exportProvider.Export(report)
	if err != nil {
		return
	}

}
