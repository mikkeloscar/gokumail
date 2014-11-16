package main

import (
	"fmt"
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/flosch/pongo2"
	"github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var tpl = loadTemplates("views")

var hashKey = securecookie.GenerateRandomKey(64)
var blockKey = securecookie.GenerateRandomKey(32)
var store = sessions.NewCookieStore(hashKey, blockKey)

func userLogin(username string, password string) error {
	return nil
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth")

	username, password := r.FormValue("username"), r.FormValue("password")

	err := userLogin(username, password)

	if err != nil {
		http.Redirect(w, r, "/", http.StatusAccepted)
		return
	}

	session.Values["user"] = username
	session.Options.MaxAge = 0 // End session when browser session ends

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/"+username, http.StatusSeeOther)
}

func settings(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth")
	user := r.URL.Query().Get(":username")

	if sess_user, ok := session.Values["user"]; ok && sess_user == user {
		// Get settings
		if r.Method == "GET" {
			settings, err := GetSettings(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			renderTemplate(w, "settings", settings)
		}

		// Save settings
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			settings := &Settings{
				user,
				r.Form.Get("workmail"),
				r.Form["from[]"],
				r.Form["to[]"],
				r.Form["blacklist[]"],
			}

			err = settings.Update()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/"+user, http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth")

	if sess_user, ok := session.Values["user"]; ok {
		if str, ok := sess_user.(string); ok {
			http.Redirect(w, r, "/"+str, http.StatusSeeOther)
			return
		}
	}

	renderTemplate(w, "index", nil)
}

// RunWebInterface runs a web interface for gokumail
func RunWebInterface(port int) {
	m := pat.New()
	m.Post("/login", http.HandlerFunc(login))
	m.Get("/:username", http.HandlerFunc(settings))
	m.Post("/:username", http.HandlerFunc(settings))
	m.Get("/", http.HandlerFunc(index))

	http.Handle("/", m)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		Log.Error("failed to start web interface: " + err.Error())
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, s *Settings) {
	if temp, ok := tpl[tmpl]; ok {
		err := temp.ExecuteWriter(pongo2.Context{"s": s}, w)
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
