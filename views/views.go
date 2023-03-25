package views

/*Holds the fetch related view handlers*/

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	// "strconv"
	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/sessions"
	"github.com/alexanderi96/leafnet/types"
)

const (
	layoutsDir   = "templates/layouts"
	templatesDir = "templates"
	extension    = "/*.gohtml"
)

var templates map[string]*template.Template

// var templates *template.Template
// var homeTemplate *template.Template
// var loginTemplate *template.Template
// var userPagetemplate *template.Template
// var peopleTemplate *template.Template
// var managePersonTemplate *template.Template
// var graphTemplate *template.Template

var c types.Context
var e error

// PopulateTemplates is used to parse all templates present in
// the templates folder
func PopulateTemplates(templatesFS embed.FS) error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	tmplFiles, err := fs.ReadDir(templatesFS, templatesDir)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(templatesFS, templatesDir+"/"+tmpl.Name(), layoutsDir+extension)
		if err != nil {
			log.Println("Error parsing template: ", err)
			return err
		}
		tmplName := strings.TrimSuffix(tmpl.Name(), filepath.Ext(tmpl.Name()))
		templates[tmplName] = pt
	}
	return nil
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
