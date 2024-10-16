package main

import (
	"fmt"
	"net/http"
	"trawix/lib"
)

func main() {
	http.HandleFunc("/events", lib.HandleEventsByNpubAndRelay)

	fmt.Println("Starting listening on :1337")
	if err := http.ListenAndServe(":1337", nil); err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
