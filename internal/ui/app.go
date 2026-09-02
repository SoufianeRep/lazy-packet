package ui

import (
	lp "lazypacket"
	"lazypacket/internal/capture"
	"lazypacket/internal/layers"
	"net"
	"strings"
	"time"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type ScreenMode int

const (
	Interfaces ScreenMode = iota
	Dashboard
)

type CaptureMsg struct {
	Handle     *capture.Handle
	FramesChan chan FrameMsg
}

type FrameMsg struct {
	Packet    *lp.Packet
	TimeStamp time.Time
}
type Entry struct {
	Packet    *lp.Packet
	TimeStamp time.Time
}
type Model struct {
	width  int
	height int

	screenMode ScreenMode

	ifaces      []net.Interface
	ifaceCursor int

	handle     *capture.Handle
	framesChan chan FrameMsg

	paused bool

	packetList viewport.Model
	detailView viewport.Model
	packets    []Entry
	selected   int
}

func InitModel() Model {
	ifaces, err := net.Interfaces()
	if err != nil {
		// do something
	}

	return Model{
		width:      0,
		height:     0,
		ifaces:     ifaces,
		packetList: newViewport(),
		detailView: newViewport(),
	}
}

func packetLines(packets []Entry) string {
	lines := make([]string, len(packets))
	for i, e := range packets {
		lines[i] = formatPacket(i+1, e, e.TimeStamp.Sub(packets[0].TimeStamp))
	}
	return strings.Join(lines, "\n")
}

func newViewport() viewport.Model {
	vp := viewport.New()
	vp.SetWidth(0)
	vp.SetHeight(0)

	return vp
}

func startCapture(ifaceName string) tea.Cmd {
	return func() tea.Msg {
		handle, err := capture.Open(ifaceName)
		if err != nil {
			return nil // TODO: handle this with global error handling
		}

		frames := make(chan FrameMsg)
		go func() {
			for {
				data, ts, err := handle.ReadFrame()
				if err != nil {
					close(frames)
					return
				}

				frames <- FrameMsg{
					Packet:    lp.NewPacket(data, &layers.Ethernet{}),
					TimeStamp: ts,
				}
			}
		}()

		return CaptureMsg{
			Handle:     handle,
			FramesChan: frames,
		}
	}
}

func handFrame(frames chan FrameMsg) tea.Cmd {
	return func() tea.Msg {
		if frame, ok := <-frames; ok {
			return frame
		}
		return nil
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.packetList.SetWidth(m.width)
		m.packetList.SetHeight(m.mainPanelHeight() / 2)

		m.detailView.SetWidth(m.width)
		m.detailView.SetHeight(m.mainPanelHeight() / 2)

		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		switch m.screenMode {
		case Dashboard:
			switch msg.String() {
			case "i":
				m.screenMode = Interfaces
			case "p":
				m.paused = !m.paused
			}
		case Interfaces:
			switch msg.String() {
			case "up", "l":
				m.ifaceCursor = (m.ifaceCursor - 1 + len(m.ifaces)) % len(m.ifaces)
			case "down", "k":
				m.ifaceCursor = (m.ifaceCursor + 1) % len(m.ifaces)
			case "enter":
				selected := m.ifaces[m.ifaceCursor]
				return m, startCapture(selected.Name)
			}
		default:
			// something here
		}

	case CaptureMsg:
		m.handle = msg.Handle
		m.framesChan = msg.FramesChan
		m.screenMode = Dashboard

		return m, handFrame(m.framesChan)

	case FrameMsg:
		if !m.paused {
			m.packets = append(m.packets, Entry{Packet: msg.Packet, TimeStamp: msg.TimeStamp})
		}

		wasAtBottom := m.packetList.AtBottom()
		m.packetList.SetContent(packetLines(m.packets))
		if wasAtBottom {
			m.packetList.GotoBottom()
		}

		return m, handFrame(m.framesChan)
	}

	return m, nil
}
