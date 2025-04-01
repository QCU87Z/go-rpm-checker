package repo

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"time"
)

type Repomd struct {
	Revision string `xml:"revision"`
	Data     []struct {
		Type     string `xml:"type,attr"`
		Checksum struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
		} `xml:"checksum"`
		Location struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
		} `xml:"location"`
	} `xml:"data"`
}

type Metadata struct {
	Packages string `xml:"packages,attr"`
	Package  []struct {
		Type     string `xml:"type,attr"`
		Name     string `xml:"name"`
		Arch     string `xml:"arch"`
		Checksum struct {
			Text  string `xml:",chardata"`
			Type  string `xml:"type,attr"`
			Pkgid string `xml:"pkgid,attr"`
		} `xml:"checksum"`
		Time struct {
			File  string `xml:"file,attr"`
			Build string `xml:"build,attr"`
		} `xml:"time"`
		Size struct {
			Package   string `xml:"package,attr"`
			Installed string `xml:"installed,attr"`
			Archive   string `xml:"archive,attr"`
		} `xml:"size"`
		Location struct {
			Href string `xml:"href,attr"`
		} `xml:"location"`
		File []struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
		} `xml:"file"`
	} `xml:"package"`
}

type Repo struct {
	Name        string
	Packages    string
	LastUpdated time.Time
	Healthly    string
}

func ProcessPrimary(url string) Metadata {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	dataValue, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(dataValue)
	gzreader, err := gzip.NewReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	output, err := io.ReadAll(gzreader)
	if err != nil {
		log.Fatal(err)
	}
	var meta Metadata
	xml.Unmarshal(output, &meta)

	return meta
}

func ProcessRepomd(url string) (string, string) {
	resp, err := http.Get(url + "repodata/repomd.xml")
	if err != nil {
		log.Fatalln(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	var repo Repomd
	xml.Unmarshal(data, &repo)
	for _, v := range repo.Data {
		if v.Type == "primary" {
			return v.Location.Href, repo.Revision
		}
	}
	return "", ""
}
