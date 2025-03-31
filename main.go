package main

import (
	"fmt"
	"go-rpm-checker/repo"
	"net/http"
	"text/template"
)

var repoPageData RepoistoryPageData
var repomds = []string{
	"repodata/repomd2.xml",
	"repodata/fedora41-everything-os.repomd.xml",
	"repodata/almalinux9-extras-os.repomd.xml",
	"repodata/almalinux9-appstream-os.repomd.xml",
}

type RepoistoryPageData struct {
	Repos []repo.Repo
}

func updateRepos() {
	for _, v := range repomds {
		fmt.Println("Processing: ", v)
		metaLocation := repo.ProcessRepomd(v)
		meta := repo.ProcessPrimary(metaLocation)
		repoPageData.Repos = append(repoPageData.Repos, repo.Repo{Name: v, Packages: meta.Packages, LastUpdated: meta.Package[0].Time.File, Healthly: true})

		// fmt.Println("Meta: ", meta.Packages)
		// PrintMemUsage()
	}
}

func main() {

	// fmt.Println(len(repoPageData.Repos))

	tmpl := template.Must(template.ParseFiles("web/table.html"))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		fmt.Println("Request: ", r.URL.Path)
		http.ServeFile(w, r, "web/favicon.ico")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Println("Request: ", r.URL.Path)
		tmpl.Execute(w, repoPageData)
	})
	http.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Println("Request: ", r.URL.Path)
		updateRepos()
		// tmpl.Execute(w, repoPageData)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	fmt.Println("Starting server on :80")
	http.ListenAndServe(":80", nil)
}
