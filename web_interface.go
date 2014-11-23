package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/mikkeloscar/goimap"
)

var templates = map[string]string{
	"index":    "index.html",
	"settings": "settings.html",
}
var tpl = loadTemplates(templates, "/usr/share/gokumail/views", "views")

var hashKey = securecookie.GenerateRandomKey(64)
var blockKey = securecookie.GenerateRandomKey(32)
var store = sessions.NewCookieStore(hashKey, blockKey)

const AUTH_COOKIE = "auth"

// authenticate user via IMAP server
func userLogin(username string, password string) error {
	service := fmt.Sprintf("%s:%d", Conf.IMAP.Server, Conf.IMAP.Port)

	conn, err := net.Dial("tcp", service)
	if err != nil {
		return err
	}

	client, err := imap.NewClient(conn, Conf.IMAP.Server)
	if err != nil {
		return err
	}

	user := fmt.Sprintf(Conf.IMAP.UsernameFmt, username)

	err = client.Login(user, password)
	if err != nil {
		return err
	}

	return client.Close() // close connection to imap server
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, AUTH_COOKIE)

	username, password := r.FormValue("username"), r.FormValue("password")

	err := userLogin(username, password)

	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		Log.Error("login error: " + err.Error())
		return
	}

	session.Values["user"] = username
	session.Options.MaxAge = 0 // End session when browser session ends

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		Log.Error("server error: " + err.Error())
		return
	}

	http.Redirect(w, r, "/"+username, http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, AUTH_COOKIE)

	session.Values = make(map[interface{}]interface{})
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		Log.Error("server error: " + err.Error())
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func settings(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, AUTH_COOKIE)
	vars := mux.Vars(r)
	user := vars["username"]

	if sessUser, ok := session.Values["user"]; ok && sessUser == user {
		// Get settings
		if r.Method == "GET" {
			settings, err := GetSettings(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				Log.Error("server error: " + err.Error())
				return
			}

			if settings == nil {
				settings = &Settings{
					User:          user,
					Workmail:      "",
					FromWhitelist: []string{},
					ToWhitelist:   []string{},
					Blacklist:     []string{},
				}

				settings.Create() // add user's settings to db
			}

			renderTemplate(w, "settings", settings)
		}

		// Save settings
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				Log.Error("server error: " + err.Error())
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
				Log.Error("server error: " + err.Error())
				return
			}

			http.Redirect(w, r, "/"+user, http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, AUTH_COOKIE)

	if sessUser, ok := session.Values["user"]; ok {
		if str, ok := sessUser.(string); ok {
			http.Redirect(w, r, "/"+str, http.StatusFound)
			return
		}
	}

	renderTemplate(w, "index", nil)
}

// RunWebInterface runs a web interface for gokumail
func RunWebInterface(port int) {
	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/logout", logout).Methods("GET")
	r.HandleFunc("/{username}", settings).Methods("GET", "POST")
	r.HandleFunc("/", index).Methods("GET")

	http.Handle("/", r)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	Log.Info("HTTP server listening on port: %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		Log.Error("failed to start web interface: " + err.Error())
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, s *Settings) {
	if temp, ok := tpl[tmpl]; ok {
		err := temp.ExecuteWriter(pongo2.Context{"s": s}, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			Log.Error("server error: " + err.Error())
		}
	} else {
		http.Error(w, "invalid template", http.StatusInternalServerError)
		Log.Error("server error: invalid template")
	}
}

func loadTemplates(templates map[string]string, dirs ...string) map[string]*pongo2.Template {
	tpl := make(map[string]*pongo2.Template)

	for _, dir := range dirs {
		for i, temp := range templates {
			path := dir + "/" + temp
			if _, err := os.Stat(path); err == nil {
				tpl[i] = pongo2.Must(pongo2.FromFile(path))
				delete(templates, i)
			}
		}
	}

	return tpl
}
