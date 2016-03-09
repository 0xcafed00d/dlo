package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func handlerEditLetter(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerEditLetter")
	http.ServeFile(w, r, "www/dloedit.html")
}

func handlerPostLetter(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerPostLetter:", r.Form["messagetext"])
	http.ServeFile(w, r, "www/dloview.html")
}

func handlerViewRandom(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerViewRandom")
	http.ServeFile(w, r, "www/dloview.html")
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

// Linux: to listen on ports <1024: sudo setcap cap_net_bind_service=+ep dlo
func main() {
	flag.Parse()

	http.HandleFunc("/", handler)

	log.Printf("Listening on port %d", config.ListenPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.ListenPort), nil)
	exitOnError(err)
}
