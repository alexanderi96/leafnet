package main

import (
	// "encoding/json"

	"embed"
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"

	//"html/template"
	"net/http"
	// "strconv"
	//"time"

	"github.com/alexanderi96/leafnet/config"
	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/sessions"
	"github.com/alexanderi96/leafnet/views"
)

var (
	//go:embed templates/*.gohtml templates/layouts/*.gohtml
	templatesFS embed.FS

	defaultConfigPath string = ".config/leafnet/config.yaml"
	configFilePath    string

	sessionKey string
)

func init() {
	flag.StringVar(&sessionKey, "session-key", "", "The Leafnet session key")
	flag.CommandLine.StringVar(&sessionKey, "s", "", "The Leafnet session key (shorthand)")

	flag.StringVar(&configFilePath, "config", "", "Path to the configuration file")
	flag.CommandLine.StringVar(&configFilePath, "c", "", "Path to the configuration file (shorthand)")
	flag.Parse()

	if sessionKey == "" {
		sessionKey = os.Getenv("LSESSION_KEY")
	}
	sessions.Init(sessionKey)

	// Check if a custom config file path is provided
	if configFilePath != "" {
		// Use the custom config file path
		config.ReadConfig(configFilePath)
	} else {
		// Use the default config file path
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		homeDir := usr.HomeDir
		if err := config.ReadConfig(filepath.Join(homeDir, defaultConfigPath)); err != nil {
			log.Fatal(err)
		}
	}

	db.Init()

	// Load the templates
	if err := views.PopulateTemplates(templatesFS); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// user relative handlers
	http.HandleFunc("/login", views.LoginFunc)
	http.HandleFunc("/signup", views.SignUpFunc)
	http.HandleFunc("/logout", views.RequiresLogin(views.LogoutFunc))
	http.HandleFunc("/my-profile", views.RequiresLogin(views.UserPage))
	http.HandleFunc("/delete-user", views.RequiresLogin(views.DeleteMyAccount))
	//http.HandleFunc("/update-user", views.RequiresLogin(views.UpdateAccountInfo))

	http.HandleFunc("/", views.RequiresLogin(views.HomeFunc))
	http.HandleFunc("/graph", views.RequiresLogin(views.GraphFunc))

	http.HandleFunc("/manage-person", views.RequiresLogin(views.AddPerson))
	http.HandleFunc("/view", views.RequiresLogin(views.ViewPeople))
	http.HandleFunc("/delete", views.RequiresLogin(views.DeletePerson))

	log.Println("Leafnet is running on port " + config.Config["leafnet_port"])
	http.ListenAndServe(":"+config.Config["leafnet_port"], nil)
}
