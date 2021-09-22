package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
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
type ViewResult struct {
	Long  string
	Short string
}

func main() {
	fmt.Println("Starting server at localhost:2212")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		short := r.URL.Path
		links := reading()

		for i := 0; i < len(links.Links); i++ {
			if links.Links[i].Shorted == short[1:] {
				fmt.Fprintf(w, "<script>document.location.href='"+links.Links[i].Origin+"'</script>")
				//fmt.Fprintf(w, links.Links[i].Origin)
				return
			}
		}

		http.ServeFile(w, r, "html/index.html")
	})
	http.HandleFunc("/createshort", func(w http.ResponseWriter, r *http.Request) {
		longUrl := r.URL.Query().Get("url")

		show := ViewResult{
			Long:  longUrl,
			Short: writing(longUrl),
		}
		tmpl, _ := template.ParseFiles("html/done.html")
		tmpl.Execute(w, show)
	})
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

	return links
}

func writing(long string) string {
	symbols := "abcdefgjhijklmnopqrstuvwxyzABCDEFGJHIJKLMNOPQRSTUVWXYZ1234567890"
	short := ""
	rand.Seed(41)
	for j := 0; j < 8; j++ {
		short += string(symbols[rand.Intn(len(symbols))])
	}
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

	return short
}
