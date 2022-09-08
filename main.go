package main

import (
	"flag"
	"fmt"
	"github.com/deathsgun/art/excel"
	"github.com/deathsgun/art/export"
	"github.com/deathsgun/art/login"
	"github.com/deathsgun/art/provider/registry"
	"github.com/deathsgun/art/redmine"
	"github.com/deathsgun/art/text"
	"github.com/deathsgun/art/untis"
	"os"
	"strings"
)

// init Register all possible providers before main is called
func init() {
	registry.ImportProviders = append(registry.ImportProviders, untis.NewUntisProvider(), redmine.NewRedmineProvider())
	registry.ExportProviders = append(registry.ExportProviders, text.NewTextProvider(), excel.NewExcelProvider())
}

func main() {
	// Setup command arguments for login and export
	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
	loginProvider := loginCmd.String("provider", "", "")
	loginUsername := loginCmd.String("username", "", "")
	loginPassword := loginCmd.String("password", "", "")

	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	exportDate := exportCmd.String("date", "", "")
	exportOutput := exportCmd.String("output", "", "")
	exportProvider := exportCmd.String("provider", "", "")
	exportPrintDates := exportCmd.Bool("print-dates", false, "")

	// Always require a sub command
	if len(os.Args) < 2 {
		printHelp()
		return
	}
	// Check if position where the subcommand is
	command := os.Args[1]
	if strings.HasPrefix(command, "-") {
		printHelp()
		return
	}

	switch command {
	case "login":
		err := loginCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		login.HandleLogin(*loginProvider, *loginUsername, *loginPassword)
		return
	case "export":
		err := exportCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		export.HandleExport(*exportProvider, *exportDate, *exportOutput, *exportPrintDates)
		return
	case "providers":
		println("Import providers:")
		for _, importProvider := range registry.ImportProviders {
			println(" - " + importProvider.Name())
		}
		println("Export providers:")
		for _, exportProvider := range registry.ExportProviders {
			println(" - " + exportProvider.Name())
		}
		return
	default:
		fmt.Printf("unknown command \"%s\" for \"art\"\n", command)
		println()
		printHelp()
		return
	}
}

func printHelp() {
	println("Usage: art <command> [flags]\n")
	println("COMMANDS")
	println("\tlogin - allows you to setup all")
	println("\texport - collects data from all import providers and exports them using the export provider")
	println("\tproviders - lists all import and export providers")
	println("COMMON FLAGS")
	println("\t--provider <name> the provider you want to use for export or login")
	println("LOGIN FLAGS")
	println("\t--username <user>")
	println("\t--password <password>")
	println("EXPORT FLAGS")
	println("\t--date (dd.MM.YYYY) the date of a week from which data should be exported")
	println("\t--output (output directory) the directory where the reports should be exported to (defaults to the current work directory)")
	println("\t--print-dates groups report entries by date and adds the date as header")
}
