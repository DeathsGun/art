package main

import (
	"flag"
	"fmt"
	"github.com/deathsgun/art/export"
	"github.com/deathsgun/art/login"
	"os"
	"strings"
)

func main() {
	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
	loginProvider := loginCmd.String("provider", "", "")
	loginUsername := loginCmd.String("username", "", "")
	loginPassword := loginCmd.String("password", "", "")

	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	exportStart := exportCmd.String("start-date", "", "")
	exportTemplate := flag.String("template", "", "")
	exportOutput := flag.String("output", "", "")
	exportProvider := flag.String("provider", "", "")

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
		export.HandleExport(*exportProvider, *exportStart, *exportTemplate, *exportOutput)
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
	println("\tlogin")
	println("\texport")
	println("COMMON FLAGS")
	println("\t--provider")
	println("LOGIN FLAGS")
	println("\t--username")
	println("\t--password")
	println("EXPORT FLAGS")
	println("\t--template")
	println("\t--output")
}
