package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/crislainesc/bookings/pkg/config"
)

var app *config.AppConfig

func NewTemplates(appConfig *config.AppConfig) {
	app = appConfig
}

// RenderTemplate renders a template using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string) {
	var tcache map[string]*template.Template
	if app.UseCache {
		tcache = app.TemplateCache
	} else {
		tcache, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := tcache[tmpl]

	if !ok {
		log.Fatal("could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	err := t.Execute(buf, nil)

	if err != nil {
		log.Println(err)
	}

	// render template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

// func RenderTemplate(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error

// 	// check to see if we already have the template in our cache
// 	_, inMap := tc[t]
// 	if !inMap {
// 		// need to create the template
// 		log.Println("creating template and adding to cache")
// 		err = createTemplateCache(t)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	} else {
// 		// we have the template in the cache
// 		log.Println("using cached template")
// 	}

// 	tmpl = tc[t]

// 	err = tmpl.Execute(w, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	// get all of the files named *.page.tmpl.html from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl.html")

	if err != nil {
		return myCache, err
	}

	// range through all of the files
	for _, page := range pages {
		// get the file name
		name := filepath.Base(page)
		// create the template
		ts, err := template.New(name).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl.html")

		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = template.ParseGlob("./templates/*.layout.tmpl.html")
			if err != nil {
				return myCache, err
			}

		}

		myCache[name] = ts
	}

	return myCache, nil
}
