package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.ServeFile(w, r, "www/dloview.html")
	} else {
		http.ServeFile(w, r, "www/dloedit.html")
	}

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
