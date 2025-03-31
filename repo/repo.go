package repo

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"io"
	"log"
	"os"
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
	LastUpdated string
	Healthly    bool
}

func ProcessPrimary(file string) Metadata {
	data, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	dataValue, err := io.ReadAll(data)
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

	// for _, v := range meta.Package {
	// 	fmt.Printf("%s %s\n", v.Location.Href, v.Checksum.Text)
	// }

	// PrintMemUsage()
	// fmt.Println("Done")
	return meta
}

func ProcessRepomd(file string) string {
	data, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	dataValue, err := io.ReadAll(data)
	if err != nil {
		log.Fatal(err)
	}
	var repo Repomd
	xml.Unmarshal(dataValue, &repo)

	for _, v := range repo.Data {
		if v.Type == "primary" {
			return (v.Location.Href)
		}
	}
	// PrintMemUsage()
	// fmt.Println("Done")

	return ""
}
