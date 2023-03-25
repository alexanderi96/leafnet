package views

/*Holds the fetch related view handlers*/

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	// "strings"
	// "strconv"
	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/sessions"
	"github.com/alexanderi96/leafnet/types"
)

// const (
// 	layoutsDir   = "public/templates"
// 	templatesDir = "public"
// 	extension    = "/*.html"
// )

// var templates map[string]*template.Template

var templates *template.Template
var homeTemplate *template.Template
var loginTemplate *template.Template
var userPagetemplate *template.Template
var peopleTemplate *template.Template
var managePersonTemplate *template.Template
var graphTemplate *template.Template

var c types.Context
var e error

// PopulateTemplates is used to parse all templates present in
// the templates folder
func PopulateTemplates() { //templatesFS embed.FS) {
	var allFiles []string
	templatesDir := "./public/templates/"
	files, err := os.ReadDir(templatesDir)
	if err != nil {
		log.Println(err)
		os.Exit(1) // No point in running app if templates aren't read
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, templatesDir+filename)
		}
	}

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	templates, err = template.ParseFiles(allFiles...)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	homeTemplate = templates.Lookup("home.html")
	loginTemplate = templates.Lookup("login.html")
	userPagetemplate = templates.Lookup("userprofile.html")
	peopleTemplate = templates.Lookup("people.html")
	managePersonTemplate = templates.Lookup("manageperson.html")
	graphTemplate = templates.Lookup("graph.html")
}

func prepareContext(w http.ResponseWriter, r *http.Request) {
	//load user info
	if c.User, e = db.GetUserInfo(sessions.GetCurrentUser(r)); e != nil {
		log.Println("Internal server error retriving user info")
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}

	//load persons
	c.Persons = db.GetPersons()
}

// func setCookie(w http.ResponseWriter) {
// 	c.CSRFToken = token

// 	cookie := http.Cookie{
// 		Name:     "csrftoken",
// 		Value:    token,
// 		Path:     "/",
// 		MaxAge:   3600,
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteLaxMode,
// 	}

// 	log.Println("Setting cookie: ", cookie)

// 	http.SetCookie(w, &cookie)
// }

// func checkCookie(r *http.Request) bool {
// 	cookie, err := r.Cookie("csrftoken")
// 	log.Println("Checking cookie: ", cookie)
// 	if err != nil {
// 		log.Println("Error getting cookie: ", err)
// 		return false
// 	}
// 	if cookie.Value != token {
// 		log.Println("Cookie is not valid")
// 		return false
// 	}
// 	if time.Now().After(cookie.Expires) {
// 		log.Println("Cookie is expired: ", time.Now(), "Cookie expiration: ", cookie.Expires)
// 		return false
// 	}

// 	return true
// }
