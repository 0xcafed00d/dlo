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

var files *FileIndex
var editTemplate *template.Template
var viewTemplate *template.Template
var random *AtomicRand

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

func handlerEditLetter(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerEditLetter")
	logOnError(editTemplate.Execute(w, ""))
}

func handlerPostLetter(w http.ResponseWriter, r *http.Request) {
	data := r.Form["messagetext"][0]
	log.Println("handlerPostLetter:", data)
	index := files.ReserveFileIndex()
	logOnError(files.StoreFile(index, data))
	logOnError(viewTemplate.Execute(w, data))
}

func handlerViewRandom(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerViewRandom")

	n := random.Int63n(files.GetFileCount())
	data, err := files.LoadFile(n)
	logOnError(err)
	logOnError(viewTemplate.Execute(w, data))
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.RequestURI)

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

func handlerFavicon(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// Linux: to listen on ports <1024: sudo setcap cap_net_bind_service=+ep dlo
func main() {
	flag.Parse()

	random = MakeAtomicRand()

	files = MakeFileIndex(config.DataFolder)
	files.RefeshFileCount()

	t, err := template.ParseFiles(
		filepath.Join(config.WWWfolder, "dloedit.html"),
		filepath.Join(config.WWWfolder, "dloview.html"),
	)
	exitOnError(err)
	log.Println(t.DefinedTemplates())

	editTemplate = t.Lookup("dloedit.html")
	viewTemplate = t.Lookup("dloview.html")

	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", handlerFavicon)

	log.Printf("Listening on port %d", config.ListenPort)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.ListenPort), nil)
	exitOnError(err)
}
