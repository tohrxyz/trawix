# Trawix
Intermediary http server that enables fetching nostr events for `npub` from `relay` without needing `nostr` tooling and/or websocket connections.

Sometimes your environment doesn't allow importing specific libs or you don't want to deal with resource intensive websockets,
for those occassions you can use this lightweight server with simple `GET` request and receive json of your notes.

### How to use
This server listens on port `:1337` and has endpoint `/events`.
It accepts 2 query params:
- npub
- relay

!!! It currently only allows to fetch `nostr.KindTimeCalendarEvent` (31923), this will be configurable later.

You can call `GET` on
```sh
curl http(s)://your_domain:1337/events?npub={some_npub}&relay={some_relay}
```

## How to run locally
1. have go installed
2. `go run main.go`
3. go to `http://localhost:1337/events?npub=<some_npub>&relay=<some_relay>`
