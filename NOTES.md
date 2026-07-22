# lazy-packet — project notes

Working notes for continuity across sessions. Update as things change; this is not
a spec, it's a snapshot of where things stand and why.

## What this is

A lazygit-style terminal dashboard (Go) for inspecting live TCP/UDP connections,
flow stats, and packet details.

## Goals for this project (why it exists)

This is a learning project, not a delivery deadline. Primary goals, in order:

1. Learn Bubble Tea and the Elm architecture (Model / Update / View) hands-on.
2. Learn how packets are actually interpreted — Ethernet/IP/TCP/UDP header parsing
   and the transport layer in general.

Implication: prefer building things step by step and understanding each piece over
having a finished implementation dropped in. Parsing logic in particular should be
hand-rolled rather than delegated to a library's built-in decoders, since that's the
point of the exercise.

## Tech stack / decisions

- Go (module `lazy-packet`, go 1.26.5)
- UI: **Bubble Tea v2** (beta) — not v1. This was an explicit correction; v2 has real
  API differences from v1, don't assume v1 docs/examples apply directly.
- Packet capture library: not yet decided (gopacket/pcap is the likely default but
  unconfirmed — needs a decision before `internal/capture` gets built).

## Planned architecture (agreed direction, not yet built)

```
cmd/lazypacket/main.go     wiring only: build a flow.Table, start tea.Program
internal/
  capture/                 I/O boundary — gets raw bytes off the wire
  protocol/                pure parsing: bytes -> Ethernet/IPv4/IPv6/TCP/UDP structs
  flow/                    aggregates decoded packets into stateful 5-tuple connections
  ui/                      Bubble Tea: Model/Update/View, depends only on flow's plain data
```

Rationale: `protocol` has no I/O and is fully unit-testable — that's where the
transport-layer learning happens. `ui` never imports `capture` or `protocol` directly,
only `flow`'s plain data, keeping the Elm-architecture loop clean.

None of this is implemented yet. `internal/ui/app.go` is currently empty
(`package ui`); `cmd/lazypacket/main.go` just prints a startup string.

## Current status

- Repo initialized, no commits yet.
- Bare skeleton only: `go.mod`, `cmd/lazypacket/main.go` (stub), `internal/ui/app.go`
  (empty).
- No dependencies added yet.
- Next step in progress: sketching the actual window/pane layout (sidebar of
  connections? main detail pane? status bar?) before writing any UI code.

## Working style for this project

- Go step by step: discuss design (layout sketches, package responsibilities, data
  flow) before scaffolding files or installing dependencies.
- Confirm exact library major versions before adding them (see the v1/v2 bubbletea
  note above) — don't assume "latest" or the most common choice.
- Let parsing/Elm-architecture logic be written collaboratively rather than handed
  over pre-built, since understanding it is the goal.

## UI Specification 
┌──────────────────────────────────────────────────────────────────────────┐
│ iface: eth0    pkt/s: 812    total: 18,402    filter: "tcp and port 443" │
│                                                         [● capturing]    │
├────────────────────────────┬─────────────────────────────────────────────┤
│ PROTOCOLS                  │ PACKET DETAILS                      #482    │
│  [x]TCP  [x]UDP  [ ]IPv6   │                                             │
├────────────────────────────┤ ▼ Ethernet II                               │
│ FILTER                     │      Src MAC   11:22:33:44:55:66            │
│  > tcp and port 443_       │      Dst MAC   aa:bb:cc:dd:ee:ff            │
├────────────────────────────┤      EtherType 0x0800 (IPv4)                │
│ PACKETS                    │                                             │
│  #    SRC:PORT    PROTO    │ ▼ IPv4                                      │
│  480  .42:51342   TCP      │      Src 192.168.1.42   Dst 93.184.216.34   │
│  481  .1:53        UDP     │      TTL 64             Proto TCP(6)        │
│ >482  .42:51342   TCP      │                                             │
│  483  .42:51344   TCP      │ ▼ TCP                                       │
│  484  .1:53        UDP     │      Src 51342   Dst 443 (https)            │
│  ...                       │      Flags [SYN]   Seq 1   Ack 1            │
│                            │                                             │
│                            │ ▶ Payload (0 bytes)                         │
│                            ├─────────────────────────────────────────────┤
│                            │ 0000  aa bb cc dd ee ff 11 22   ..."iVf..   │
│                            │ 000e  45 00 00 3c 1c 46 40 00   E..<.F@.    │
├────────────────────────────┴─────────────────────────────────────────────┤
│ ↑/↓ move   enter select   / filter   p pause   c clear   q quit          │
└──────────────────────────────────────────────────────────────────────────┘