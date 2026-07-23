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

## Current status (2026-07-24)

- `internal/ui`: basic Model/Update/View skeleton and main program wiring exist
  (bubbletea v2), committed.
- `internal/protocol/ethernet.go`: in progress.
  - Types defined: `EthernetFrame` (DstMAC, SrcMAC, VLAN *VLANTag, EtherType),
    `VLANTag` (TPID + embedded `TCI`), `TCI` (PCP, DEI, VID) — handles optional
    802.1Q tagging between src MAC and EtherType.
  - `ParseEthernetFrame(data []byte) (*EthernetFrame, []byte, error)` — signature
    decided, body not yet implemented (still a stub returning a zero struct).
  - `sample.pcap` captured and committed (own traffic) to use as a parsing fixture
    source — individual packet fixtures not yet extracted as Go `[]byte` literals.
- Open/deferred decisions (revisit later, not blocking current work):
  - `net.HardwareAddr` vs. plain `[6]byte` for MAC fields — currently `[6]byte`,
    considered but not changed.
  - A `Contents []byte` (raw header bytes) field per layer, similar to gopacket's
    `BaseLayer` — wanted for a future "click a layer in the tree, highlight its
    bytes in the hex dump" UI feature. Not added yet; not needed until the UI
    hex-dump work starts.
  - Method-receiver + interface-based decode architecture (gopacket's
    `DecodeFromBytes` / `DecodingLayer` pattern) — deliberately deferred. Current
    free-function style (`ParseX(data) -> (*X, []byte, error)`) is simpler to
    test and is fine for now; converting later is a low-risk mechanical refactor
    since the actual byte-parsing logic doesn't change, only the packaging. The
    full interface-based polymorphic dispatch is its own separate future learning
    goal, not part of this pass.
  - Packet capture library still undecided (gopacket/pcap likely, unconfirmed).

## Where to pick up next time

Working through the todo list for the "hand-rolled protocol parsing" phase, in order:

1. Extract 2-3 real packets from `sample.pcap` as raw hex byte fixtures (one TCP,
   one UDP; cross-check each against Wireshark's decode as ground truth).
2. Finish `ParseEthernetFrame`'s body (currently a stub) — including the VLAN-tag
   peek/skip logic.
3. IPv4 header parsing (`internal/protocol/ipv4.go`).
4. TCP header parsing (`internal/protocol/tcp.go`).
5. UDP header parsing (`internal/protocol/udp.go`).
6. Unit tests for all of the above against the fixtures from step 1.
7. A throwaway harness that runs one fixture through the full Ethernet → IPv4 →
   TCP/UDP chain and prints every field, checked by eye against Wireshark. This is
   the finish line for this phase — `internal/capture` and the UI don't get
   touched until this works end to end.

## Working style for this project

- Go step by step: discuss design (layout sketches, package responsibilities, data
  flow) before scaffolding files or installing dependencies.
- Confirm exact library major versions before adding them (see the v1/v2 bubbletea
  note above) — don't assume "latest" or the most common choice.
- The user writes all code themselves (parsing logic, Elm-architecture logic,
  everything). Explicit correction (2026-07-24): "I don't want you to code
  anything, I want to write the code myself." Claude's role is explaining
  concepts, answering design questions, and reviewing/commenting code on request
  — not authoring it, unless a patch is explicitly requested (e.g. adding doc
  comments to an existing file).

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