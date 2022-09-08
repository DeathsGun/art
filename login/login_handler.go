package login

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deathsgun/art/provider"
	"github.com/deathsgun/art/provider/registry"
	"golang.org/x/term"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

var logins []*Login
var changed = false

func init() {
	config, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Failed to get user config dir: %v", err)
		os.Exit(1)
		return
	}
	f, err := os.OpenFile(filepath.Join(config, ".art.json"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logins = []*Login{}
			saveLogins()
			return
		}
		fmt.Printf("Failed to create config file: %v", err)
		os.Exit(1)
		return
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	err = json.NewDecoder(f).Decode(&logins)
	if err != nil {
		fmt.Printf("Failed to decode config file: %v", err)
		os.Exit(1)
		return
	}
}

func HandleLogin(prov string, username string, password string) {
	if prov == "" {
		if username != "" || password != "" {
			println("can't accept username and password without an specific provider")
			os.Exit(1)
		}
		handleLoginForAll()
		if changed {
			println("Successfully saved the new credentials")
		} else {
			println("Nothing changed but you can set a password for the provider like this:")
			println("art login --provider <provider> --username <user> --password <password>")
		}
		return
	}
	var p provider.Provider = nil
	for _, ip := range registry.ImportProviders {
		if ip.Name() == prov {
			p = ip
		}
	}
	if p == nil {
		for _, exportProvider := range registry.ExportProviders {
			if exportProvider.Name() == prov {
				p = exportProvider
			}
		}
	}
	if p == nil {
		fmt.Printf("Provider %s not found\n", prov)
		os.Exit(1)
		return
	}

	if username == "" && password == "" {
		handleLogin(p, true)
		return
	}
	username, password, err := p.ValidateLogin(username, password)
	if err != nil {
		fmt.Printf("Error while logging in with provided credentials: %v\n", err)
		return
	}
	setCredentials(p.Name(), username, password)
	saveLogins()
	println("Successfully logged in with the provided credentials\n")
}

func handleLoginForAll() {
	for _, prov := range registry.ImportProviders {
		handleLogin(prov, false)
	}
	for _, exportProvider := range registry.ExportProviders {
		handleLogin(exportProvider, false)
	}
	saveLogins()
}

func handleLogin(prov provider.Provider, force bool) {
	if !prov.NeedsLogin() {
		return
	}
	login := getLogin(prov.Name())
	if login != nil && !force {
		_, _, err := prov.ValidateLogin(login.Username, login.Password)
		if err == nil {
			return
		}
		fmt.Printf("Login for %s is invalid\n", prov.Name())
	} else if login == nil {
		fmt.Printf("Login for %s not found\n", prov.Name())
	}

	if !force {
		fmt.Printf("Do you wan't to configure %s? [y/N] ", prov.Name())
		var yesNo string
		_, _ = fmt.Scanf("%s\n", &yesNo)
		if strings.ToLower(yesNo) != "y" {
			fmt.Printf("Oke skipping configuration for %s\n", prov.Name())
			saveLogins()
			return
		}
	}

	for i := 0; i < 3; i++ {
		username := requireInput("Username: ")
		password := requirePassword("Password: ")
		println()
		username, password, err := prov.ValidateLogin(username, password)
		if err == nil {
			println("Successfully logged in with the provided credentials")
			setCredentials(prov.Name(), username, password)
			saveLogins()
			return
		}
		fmt.Printf("\nError while logging in with provided credentials: %v\n", err)
	}
	println("3 incorrect attempts skipping.")
}

func GetLogin(providerName string) (string, string) {
	for _, login := range logins {
		if login.Name == providerName {
			return login.Username, login.Password
		}
	}

	return "", ""
}

func setCredentials(name string, username string, password string) {
	changed = true
	for _, l := range logins {
		if l.Name == name {
			l.Username = username
			l.Password = password
			return
		}
	}
	logins = append(logins, &Login{
		Name:     name,
		Username: username,
		Password: password,
	})

}

func requireInput(text string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(text)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(strings.ReplaceAll(input, "\n", ""), "\r", "")
}

func requirePassword(text string) string {
	fmt.Printf(text)
	//goland:noinspection GoRedundantConversion
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return ""
	}
	return string(password)
}

func getLogin(name string) *Login {
	for _, login := range logins {
		if login.Name == name {
			return login
		}
	}
	return nil
}

func saveLogins() {
	config, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Failed to get config dir: %v", err)
		os.Exit(1)
		return
	}
	f, err := os.OpenFile(filepath.Join(config, ".art.json"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to open config file: %v", err)
		os.Exit(1)
		return
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	err = json.NewEncoder(f).Encode(logins)
	if err != nil {
		fmt.Printf("Failed to write config file: %v", err)
		os.Exit(1)
		return
	}
}
