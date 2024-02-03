package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	c := getCookie(w, r)
	if r.Method == http.MethodPost {
		f, fh, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		ext := strings.Split(fh.Filename, ".")[1]
		h := sha1.New()
		io.Copy(h, f)
		fn := fmt.Sprintf("%x", h.Sum(nil)) + "." + ext

		wd, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}

		fp := filepath.Join(wd, "public", "pics", fn)
		nf, err := os.Create(fp)
		if err != nil {
			log.Fatalln(err)
		}
		defer nf.Close()

		f.Seek(0, 0)
		io.Copy(nf, f)
		c = appendValue(w, c, fn)
	}

	xs := strings.Split(c.Value, "|")
	tpl.ExecuteTemplate(w, "index.gohtml", xs)
}

func appendValue(w http.ResponseWriter, c *http.Cookie, fn string) *http.Cookie {
	s := c.Value
	if !strings.Contains(s, fn) {
		s += "|" + fn
	}
	c.Value = s
	http.SetCookie(w, c)
	return c
}

func getCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	c, err := r.Cookie("session")
	if err != nil {
		v := uuid.NewString()
		c = &http.Cookie{
			Name:     "session",
			Value:    v,
			HttpOnly: true,
		}
		http.SetCookie(w, c)
	}
	return c
}
