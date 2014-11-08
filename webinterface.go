package main

import (
	"fmt"
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/flosch/pongo2"
)

var tpl = loadTemplates("views")

func login(w http.ResponseWriter, r *http.Request) {

}

func settings(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Query().Get(":username"))
	renderTemplate(w, "settings", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

// RunWebInterface runs a web interface for gokumail
func RunWebInterface(port int) {
	m := pat.New()
	m.Post("/login", http.HandlerFunc(login))
	m.Get("/:username", http.HandlerFunc(settings))
	m.Get("/", http.HandlerFunc(index))

	http.Handle("/", m)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func renderTemplate(w http.ResponseWriter, tmpl string, s *Settings) {
	if temp, ok := tpl[tmpl]; ok {
		err := temp.ExecuteWriter(nil, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "invalid template", http.StatusInternalServerError)
	}
}

func loadTemplates(dir string) map[string]*pongo2.Template {
	tpl := make(map[string]*pongo2.Template)
	tpl["index"] = pongo2.Must(pongo2.FromFile(dir + "/index.html"))
	tpl["settings"] = pongo2.Must(pongo2.FromFile(dir + "/settings.html"))
	return tpl
}
