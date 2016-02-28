package main

import "net/http"

func editCommandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.ServeFile(w, r, "www/dloview.html")
	} else {
		http.ServeFile(w, r, "www/dloedit.html")
	}

}

func main() {
	http.HandleFunc("/", editCommandHandler)
	http.ListenAndServe(":8080", nil)
}
