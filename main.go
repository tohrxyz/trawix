package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func fetchNotesOnDemand(npub string, relayUrl string) ([]string, error) {
	ctx := context.Background()

	relay, err := nostr.RelayConnect(ctx, relayUrl)
	if err != nil {
		return nil, err
	}
	var filters nostr.Filters
	if _, v, err := nip19.Decode(npub); err == nil {
		pub := v.(string)
		filters = []nostr.Filter{{
			Kinds:   []int{nostr.KindTimeCalendarEvent},
			Authors: []string{pub},
		}}
	} else {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sub, err := relay.Subscribe(ctx, filters)
	if err != nil {
		return nil, err
	}

	var events []string
	for ev := range sub.Events {
		json, err := json.MarshalIndent(ev, "", " ")
		if err != nil {
			fmt.Printf("Error with parsing event: %v", err)
		}
		events = append(events, string(json))
	}
	return events, nil
}

func handleEventsByNpubAndRelay(w http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	queries := r.URL.Query()
	npub := queries.Get("npub")
	if npub == "" {
		fmt.Printf("No npub found.")
		http.Error(w, "No npub found", http.StatusBadRequest)
		return
	}
	relay := queries.Get("relay")
	if relay == "" {
		fmt.Printf("No relay found.")
		http.Error(w, "No relay found", http.StatusBadRequest)
		return
	}

	events, err := fetchNotesOnDemand(npub, relay)
	if err != nil {
		fmt.Printf("Error while fetching notes: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data, err := json.Marshal(events)
	if err != nil {
		fmt.Printf("Error serializing json: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	timeEnd := time.Now()
	diff := timeEnd.UnixMilli() - timeStart.UnixMilli()
	sliceNpub := npub[:8] + "..." + npub[len(npub)-8:]
	fmt.Printf("Success for %s from %s -> took %vms\n", sliceNpub, relay, diff)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

func main() {
	http.HandleFunc("/events", handleEventsByNpubAndRelay)

	fmt.Println("Starting listening on :1337")
	if err := http.ListenAndServe(":1337", nil); err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
