package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Links struct {
	XMLName xml.Name `xml:"links"`
	Links   []Link   `xml:"link"`
}
type Link struct {
	XMLName xml.Name `xml:"link"`
	Origin  string   `xml:"origin"`
	Shorted string   `xml:"shorted"`
}

func main() {
	fmt.Println("Starting server at localhost:2212")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/index.html")
	})
	http.HandleFunc("/createshort", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html/done.html")
		longUrl := r.URL.Query().Get("url")
		fmt.Println("url: " + longUrl)
	})
	links := reading()
	for i := 0; i < len(links.Links); i++ {
		http.HandleFunc(links.Links[i].Shorted, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, links.Links[i].Origin)
		})
	}
	fmt.Println("Server is listening. . .")
	http.ListenAndServe(":2212", nil)
}

func reading() Links {
	xmFile, err := os.Open("data.xml")

	if err != nil {
		fmt.Println(err)
	}

	defer xmFile.Close()

	byteValue, _ := ioutil.ReadAll(xmFile)

	var links Links

	xml.Unmarshal(byteValue, &links)

	for i := 0; i < len(links.Links); i++ {
		fmt.Println("origin: " + links.Links[i].Origin)
		fmt.Println("shorted: " + links.Links[i].Shorted)
	}

	return links
}

func writing(long string, short string) {
	xmFile, err := os.Open("data.xml")

	if err != nil {
		fmt.Println(err)
	}

	defer xmFile.Close()

	byteValue, _ := ioutil.ReadAll(xmFile)

	var links Links

	xml.Unmarshal(byteValue, &links)

	link := Link{
		Origin:  long,
		Shorted: short,
	}

	links.Links = append(links.Links, link)

	file, _ := xml.MarshalIndent(links, "", "")
	_ = ioutil.WriteFile("data.xml", file, 0644)
}
