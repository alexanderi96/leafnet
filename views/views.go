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

var c types.Context

// PopulateTemplates is used to parse all templates present in
// the templates folder
func PopulateTemplates(templatesFS embed.FS) error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	tmplFiles, err := fs.ReadDir(templatesFS, templatesDir)
	if err != nil {
		return err
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(templatesFS, templatesDir+"/"+tmpl.Name(), layoutsDir+extension)
		if err != nil {
			return err
		}
		tmplName := strings.TrimSuffix(tmpl.Name(), filepath.Ext(tmpl.Name()))
		templates[tmplName] = pt
	}
	return nil
}

func prepareContext(w http.ResponseWriter, r *http.Request) {
	var err error
	if user := sessions.GetCurrentUser(r); user != "" {
		if c.User, err = db.GetUserInfoByEmail(user); err != nil {
			WriteError(w, err)
			return

		} else if c.Persons, err = db.GetPersons(); err != nil {
			WriteError(w, err)
			return

		}
	} else {
		log.Println("No user logged in, emptying context!")
		c = types.Context{}
	}
}

// write a function to handle errors
func WriteError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	log.Println(err)
	c.Error = err
	if err := templates["error"].Execute(w, c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
