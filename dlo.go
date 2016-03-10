package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var EditTemplate *template.Template
var ViewTemplate *template.Template

func handlerEditLetter(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerEditLetter")
	logOnError(EditTemplate.Execute(w, ""))
}

func handlerPostLetter(w http.ResponseWriter, r *http.Request) {
	data := r.Form["messagetext"][0]
	log.Println("handlerPostLetter:", data)
	logOnError(ViewTemplate.Execute(w, data))
}

func handlerViewRandom(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerViewRandom")
	data := "this is some text"
	logOnError(ViewTemplate.Execute(w, data))
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.Method == "POST" {
		if _, ok := r.Form["postletter"]; ok {
			handlerPostLetter(w, r)
			return
		}

		if _, ok := r.Form["viewrandom"]; ok {
			handlerViewRandom(w, r)
			return
		}
	}
	handlerEditLetter(w, r)
}

type Config struct {
	ListenPort int
	DataFolder string
	WWWfolder  string
}

var config Config

func init() {
	flag.IntVar(&config.ListenPort, "p", 80, "http listen port")
	flag.StringVar(&config.DataFolder, "d", "./data", "data folder location")
	flag.StringVar(&config.WWWfolder, "w", "./www", "www folder location")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: dlo [options]")
		flag.PrintDefaults()
	}
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func logOnError(err error) {
	if err != nil {
		log.Println(err)
	}
}

// Linux: to listen on ports <1024: sudo setcap cap_net_bind_service=+ep dlo
func main() {
	flag.Parse()

	t, err := template.ParseFiles(
		filepath.Join(config.WWWfolder, "dloedit.html"),
		filepath.Join(config.WWWfolder, "dloview.html"),
	)
	exitOnError(err)
	log.Println(t.DefinedTemplates())

	EditTemplate = t.Lookup("dloedit.html")
	ViewTemplate = t.Lookup("dloview.html")

	http.HandleFunc("/", handler)

	log.Printf("Listening on port %d", config.ListenPort)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.ListenPort), nil)
	exitOnError(err)
}
