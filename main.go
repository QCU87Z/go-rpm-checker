package main

import (
	"fmt"
	"go-rpm-checker/repo"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var repoPageData RepoistoryPageData
var repoUrls = []string{
	"https://repo.nagios.com/nagios/9/",
	"https://repo.almalinux.org/almalinux/9/extras/x86_64/os/",
	// "https://download.rockylinux.org/pub/rocky/9/extras/x86_64/os/",
	// "https://repo.almalinux.org/almalinux/9/BaseOS/x86_64/os/",
	// "https://download.rockylinux.org/pub/rocky/8/extras/x86_64/os/",
	// "https://repo.almalinux.org/almalinux/8/BaseOS/x86_64/os/",
	// "https://dl.fedoraproject.org/pub/epel/9/Everything/x86_64/",
}

type RepoistoryPageData struct {
	Repos []repo.Repo
}

func updateRepos() {
	repoPageData.Repos = []repo.Repo{}
	for _, v := range repoUrls {
		fmt.Println("Processing: ", v)
		metaLocation, revision := repo.ProcessRepomd(v)
		meta := repo.ProcessPrimary(v + metaLocation)
		i, err := strconv.ParseInt(revision, 10, 64)
		if err != nil {
			panic(err)
		}
		tm := time.Unix(i, 0)
		weekAgo := time.Now().Add(-168 * time.Hour)
		var health string
		var bad bool
		if !tm.After(weekAgo) {
			health = "‼️"
			bad = true
		} else {
			health = "✅"
			bad = false
		}
		repoPageData.Repos = append(repoPageData.Repos, repo.Repo{Name: v, Packages: meta.Packages, LastUpdated: tm, Healthly: health, Bad: bad})
	}
}

func main() {
	tmpl := template.Must(template.ParseFiles("web/table.html"))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		fmt.Println("Request: ", r.URL.Path)
		http.ServeFile(w, r, "web/favicon.ico")
	})
	http.HandleFunc("/output.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		fmt.Println("Request: ", r.URL.Path)
		http.ServeFile(w, r, "web/output.css")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Println("Request: ", r.URL.Path)
		tmpl.Execute(w, repoPageData)
	})
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Println("Request: ", r.URL.Path)
		updateRepos()
		// tmpl.Execute(w, repoPageData)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	fmt.Println("Starting server on :80")
	http.ListenAndServe(":80", nil)
}
