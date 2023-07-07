package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/crislainesc/bookings/internal/config"
	"github.com/crislainesc/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"formatDate":           FormatDate,
	"formatDateWithLayout": FormatDateWithLayout,
	"iterate":              Iterate,
	"add":                  Add,
}

var app *config.AppConfig
var templatesPath = "../../templates/"

// NewRenderer  sets the config for the template package
func NewRenderer(appConfig *config.AppConfig) {
	app = appConfig
}

func FormatDate(d time.Time) string {
	return d.Format("2006-01-02")
}

func Add(a, b int) int {
	return a + b
}

// Iterate returns a slice of integers, starting at 0, going to count
func Iterate(count int) []int {
	var items []int
	for i := 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

func FormatDateWithLayout(t time.Time, l string) string {
	return t.Format(l)
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

// Template  renders a template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tcache map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tcache = app.TemplateCache
	} else {
		tcache, _ = CreateTemplateCache()
	}

	t, ok := tcache[tmpl]
	if !ok {
		//log.Fatal("could not get template from template cache")
		return errors.New("could not get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}

	return nil
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(templatesPath + "*.page.tmpl.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(templatesPath + "*.layout.tmpl.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(templatesPath + "*.layout.tmpl.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
