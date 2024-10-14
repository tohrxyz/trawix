package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

var events []string

func fetchNotes(globalEvents *[]string, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx := context.Background()

	relayUrl := "wss://relay.damus.io"
	relay, err := nostr.RelayConnect(ctx, relayUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Relay url: ", relay.URL)

	npub := "npub"
	fmt.Println("Npub: ", npub)
	var filters nostr.Filters
	if _, v, err := nip19.Decode(npub); err == nil {
		pub := v.(string)
		filters = []nostr.Filter{{
			Kinds:   []int{nostr.KindTimeCalendarEvent},
			Authors: []string{pub},
		}}
	} else {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sub, err := relay.Subscribe(ctx, filters)
	if err != nil {
		panic(err)
	}

	var events []string
	for ev := range sub.Events {
		json, err := json.MarshalIndent(ev, "", " ")
		if err != nil {
			fmt.Printf("Error with parsing event: %v", err)
		}
		events = append(events, string(json))
	}
	*globalEvents = append(*globalEvents, events...)
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(events)
	if err != nil {
		fmt.Printf("Error serializing json: %v", err)
	}
	w.Write([]byte(data))
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go fetchNotes(&events, &wg)

	wg.Wait()
	fmt.Println(events)

	http.HandleFunc("/events", handleEvents)
	http.ListenAndServe(":1337", nil)
}
