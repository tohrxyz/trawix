# Trawix
Intermediary http server that enables fetching nostr events for `npub` from `relay` without needing `nostr` tooling and/or websocket connections.

Sometimes your environment doesn't allow importing specific libs or you don't want to deal with resource intensive websockets,
for those occassions you can use this lightweight server with simple `GET` request and receive json of your notes.

### How to use
This server listens on port `:1337` and has endpoint `/events`.
It accepts query params:
- npub (npub....)
- relay (relay.domain.com)
- kinds (0,5,31923...)

You can call `GET` on
```sh
curl http(s)://your_domain:1337/events?npub={some_npub}&relay={some_relay}&kinds={number,number,number...}
```

## How to run locally
1. have go installed
2. `go run main.go`
3. go to `http://localhost:1337/events?npub=<some_npub>&relay=<some_relay>&kinds=<number,number...>`
