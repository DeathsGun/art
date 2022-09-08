package main

import (
	"flag"
	"fmt"
	"github.com/deathsgun/art/export"
	"github.com/deathsgun/art/login"
	"github.com/deathsgun/art/provider"
	"github.com/deathsgun/art/provider/registry"
	"github.com/deathsgun/art/redmine"
	"github.com/deathsgun/art/untis"
	"os"
	"strings"
)

func init() {
	registry.ImportProviders = append(registry.ImportProviders, untis.NewUntisProvider(), redmine.NewRedmineProvider())
	registry.ExportProviders = append(registry.ExportProviders, provider.NewTextProvider(), provider.NewExcelProvider())
}

func main() {
	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
	loginProvider := loginCmd.String("provider", "", "")
	loginUsername := loginCmd.String("username", "", "")
	loginPassword := loginCmd.String("password", "", "")

	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	exportStart := exportCmd.String("start-date", "", "")
	exportOutput := exportCmd.String("output", "", "")
	exportProvider := exportCmd.String("provider", "", "")

	if len(os.Args) < 2 {
		printHelp()
		return
	}
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
		export.HandleExport(*exportProvider, *exportStart, *exportOutput)
		return
	case "providers":
		println("Import providers:")
		for _, importProvider := range registry.ImportProviders {
			println("\t" + importProvider.Name())
		}
		println("Export providers:")
		for _, exportProvider := range registry.ExportProviders {
			println("\t" + exportProvider.Name())
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
	println("\t--start-date (dd.MM.YYYY) the date of a week from which data should be exported")
	println("\t--output (output directory) the directory where the reports should be exported to (defaults to the current work directory)")
}
