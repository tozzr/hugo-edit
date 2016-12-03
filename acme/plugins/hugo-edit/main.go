package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type Page struct {
	Path  string
	Title string
	Body  []byte
}

func printFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Print(err)
		return nil
	}
	fmt.Println(path)
	return nil
}

func pageListHandler(w http.ResponseWriter, r *http.Request) {
	pages := []Page{{Title: "hello", Path: "/post/hello"}}
	err := filepath.Walk("../../content", printFile)
	if err != nil {
		log.Fatal(err)
	}
	renderTemplate(w, "list", pages)
}

func pageEditHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("p")
	p, err := loadPage(path)
	if err != nil {
		p = &Page{Title: "not found"}
	}
	renderTemplate(w, "edit", p)
}

func pageSaveHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("p")
	body := r.FormValue("body")
	p := &Page{Title: path, Path: path, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/", http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t := template.Must(template.ParseFiles("tmpl/head.html", "tmpl/foot.html", "tmpl/list.html", "tmpl/edit.html"))
	t.ExecuteTemplate(w, tmpl, data)
}

func (p *Page) save() error {
	filename := "../../content" + p.Path + ".md"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(path string) (*Page, error) {
	filename := "../../content" + path + ".md"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: path + ".md", Path: path, Body: body}, nil
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", pageListHandler).Methods("GET")
	r.HandleFunc("/page", pageEditHandler).Methods("GET")
	r.HandleFunc("/page", pageSaveHandler).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", r)
	http.ListenAndServe(":1314", nil)
}
