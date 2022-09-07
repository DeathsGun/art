package export

import (
	"fmt"
	"github.com/deathsgun/art/provider"
	"github.com/deathsgun/art/provider/registry"
	"os"
	"sort"
	"time"
)

func HandleExport(prov string, start string, output string) {
	if prov == "" {
		println("a provider for export is required")
		os.Exit(1)
	}
	var err error
	t := time.Now()
	if start != "" {
		t, err = time.Parse("02.01.2006", start)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
	}

	var exportProvider provider.ExportProvider = nil
	for _, p := range registry.ExportProviders {
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

	for _, iprov := range registry.ImportProviders {
		entries, err := iprov.Import(t)
		if err != nil {
			fmt.Printf("Skipping import provider %s because it errored: %v\n", iprov.Name(), err)
			continue
		}
		report.Entries = append(report.Entries, entries...)
	}

	if len(report.Entries) == 0 {
		println("Skipping export because no provider returned data")
		os.Exit(0)
	}

	sort.Slice(report.Entries, func(i, j int) bool {
		return report.Entries[i].Date.Before(report.Entries[j].Date)
	})

	err = exportProvider.Export(report, t, output)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	println("Successfully exported data")
}
