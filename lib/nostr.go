package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func FetchNotesOnDemand(npub string, relayUrl string, kinds []int) ([]string, error) {
	ctx := context.Background()

	relay, err := nostr.RelayConnect(ctx, relayUrl)
	if err != nil {
		return nil, err
	}
	var filters nostr.Filters
	if _, v, err := nip19.Decode(npub); err == nil {
		pub := v.(string)
		filters = []nostr.Filter{{
			Kinds:   kinds,
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
