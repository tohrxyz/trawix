package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HandleEventsByNpubAndRelay(w http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	queries := r.URL.Query()
	npub := queries.Get("npub")
	if npub == "" {
		fmt.Printf("No npub found.\n")
		http.Error(w, "No npub found", http.StatusBadRequest)
		return
	}

	relay := queries.Get("relay")
	if relay == "" {
		fmt.Printf("No relay found.\n")
		http.Error(w, "No relay found", http.StatusBadRequest)
		return
	}

	kinds := strings.Split(queries.Get("kinds"), ",")
	if len(kinds) == 0 {
		fmt.Printf("No event kinds found.\n")
		http.Error(w, "No kinds found", http.StatusBadRequest)
		return
	}

	var parsedKinds []int
	for i, v := range kinds {
		val, err := strconv.ParseInt(kinds[i], 10, 64)
		if err != nil {
			fmt.Printf("Unable to parse %s kind\n", v)
			continue
		}
		parsedKinds = append(parsedKinds, int(val))
	}

	if len(parsedKinds) == 0 {
		fmt.Printf("No usable kinds for events.\n")
		http.Error(w, "No usable kinds for events", http.StatusBadRequest)
		return
	}

	events, err := FetchNotesOnDemand(npub, relay, parsedKinds)
	if err != nil {
		fmt.Printf("Error while fetching notes: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(events)
	if err != nil {
		fmt.Printf("Error serializing json: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timeEnd := time.Now()
	diff := timeEnd.UnixMilli() - timeStart.UnixMilli()
	sliceNpub := npub[:8] + "..." + npub[len(npub)-8:]
	fmt.Printf("Success for %s from %s -> took %vms\n", sliceNpub, relay, diff)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}
