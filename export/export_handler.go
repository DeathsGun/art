package export

import (
	"errors"
	"fmt"
	"github.com/deathsgun/art/provider"
	"github.com/deathsgun/art/provider/registry"
	"github.com/deathsgun/art/utils"
	"os"
	"path/filepath"
	"time"
)

// HandleExport contains the control logic for the export
//
// Params:
//
//	prov is the provider requested by the user
//	date is some in date in a week which will be used to export data
//	output is the directory where export files are saved
func HandleExport(prov string, date string, output string, printDates bool) {
	if prov == "" {
		println("a provider for export is required")
		os.Exit(1)
	}
	var err error
	dateTime := time.Now()
	if date != "" {
		dateTime, err = time.Parse("02.01.2006", date)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
	}

	dateTime = utils.LerpToPreviousMonday(dateTime)
	fmt.Printf("Using %s as start date\n", dateTime.Format("Monday 02.01.2006"))

	if output == "" {
		output, _ = filepath.Abs(".")
	}

	var exportProvider provider.ExportProvider = nil
	for _, p := range registry.ExportProviders {
		if p.Name() == prov {
			exportProvider = p
		}
	}

	if exportProvider == nil {
		fmt.Printf("provider %s not found\n", prov)
		os.Exit(1)
	}

	report := &provider.Report{
		Id:      0,
		Entries: []*provider.Entry{},
	}

	for _, iprov := range registry.ImportProviders {
		entries, err := iprov.Import(dateTime)
		if err != nil {
			if errors.Is(err, provider.ErrNoLoginConfigured) {
				fmt.Printf("[%s]: Skipping import. Reason: not configured\n", iprov.Name())
				continue
			}
			fmt.Printf("Skipping import provider %s because it errored: %v\n", iprov.Name(), err)
			continue
		}
		if len(entries) == 0 {
			fmt.Printf("[%s]: Skipping import. Reason: no entries provided\n", iprov.Name())
			continue
		}
		report.Entries = append(report.Entries, entries...)
	}

	if len(report.Entries) == 0 {
		println("Skipping export because no provider returned data")
		os.Exit(0)
	}

	err = exportProvider.Export(report, dateTime, output, printDates)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	println("Successfully exported data")
}
