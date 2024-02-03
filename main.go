package main

import (
	"github.com/google/uuid"
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

type user struct {
	UserName string
	Email    string
}

var dbSessions = make(map[string]user)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	u, ok := dbSessions[c.Value]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "index.gohtml", u)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		u := user{
			UserName: r.FormValue("username"),
			Email:    r.FormValue("email"),
		}
		v := uuid.NewString()
		http.SetCookie(w, &http.Cookie{Name: "session", Value: v, HttpOnly: true})
		dbSessions[v] = u
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	_, err := r.Cookie("session")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}
