package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

var (
	calendarEvents []string
	mu             sync.Mutex
)

func fetchNotes(npub string) {
	ctx := context.Background()

	relayUrl := "wss://relay.damus.io"
	relay, err := nostr.RelayConnect(ctx, relayUrl)
	if err != nil {
		panic(err)
	}
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
	mu.Lock()
	calendarEvents = events
	mu.Unlock()
}

func startPeriodicFetch(interval time.Duration, npub string) {
	timezone := time.FixedZone("GMT+2", 2*60*60)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		currentTime := time.Now().In(timezone)
		fmt.Printf("Refreshing events cache @ %s\n", currentTime.Format("2006-01-02 15:04:05"))
		<-ticker.C
		fetchNotes(npub)
	}
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(calendarEvents)
	if err != nil {
		fmt.Printf("Error serializing json: %v", err)
	}
	w.Write([]byte(data))
}

func main() {
	npub := os.Args[1]
	go fetchNotes(npub)
	go startPeriodicFetch(60*time.Second, npub)

	http.HandleFunc("/events", handleEvents)

	fmt.Println("Starting listening on :1337")
	if err := http.ListenAndServe(":1337", nil); err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
