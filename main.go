package main

import (
	"flag"
	"fmt"
	"github.com/deathsgun/art/login"
	"os"
	"strings"
)

var (
	providerFlag = flag.String("provider", "", "")
	usernameFlag = flag.String("username", "", "")
	passwordFlag = flag.String("password", "", "")
	templateFlag = flag.String("template", "", "")
	outputFlag   = flag.String("output", "", "")
)

func init() {
	flag.Usage = printHelp
	flag.Parse()
}

func main() {
	if len(os.Args) == 1 {
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
		login.HandleLogin(*providerFlag, *usernameFlag, *passwordFlag)
		return
	case "export":
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
