package main

import (
	"html/template"
	"net/http"
)

type Repoistory struct {
	Repo        string
	Release     string
	LastUpdated int
	Healthly    bool
}

type RepoistoryPageData struct {
	Repos []Repoistory
}

func main() {
	data := RepoistoryPageData{
		Repos: []Repoistory{
			{Repo: "AlmaLinux", Release: "9", LastUpdated: 17232324, Healthly: true},
			{Repo: "AlmaLinux", Release: "9.1", LastUpdated: 17232320, Healthly: true},
			{Repo: "Fedora", Release: "41", LastUpdated: 17232220, Healthly: true},
		},
	}

	tmpl := template.Must(template.ParseFiles("table.html"))
	// tmpl := template.Must(template.ParseFiles("accord.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, data)
	})
	http.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		a := Repoistory{Repo: "abc", Release: "1", LastUpdated: 12345, Healthly: false}
		data.Repos = append(data.Repos, a)
		tmpl.Execute(w, data)
	})
	http.ListenAndServe(":80", nil)
}
