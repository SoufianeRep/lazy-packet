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

type FocusedPanel int

const (
	NoFocus FocusedPanel = iota
	PacketsFocus
	DetailsFocus
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

	handle *capture.Handle

	framesChan chan FrameMsg
	paused     bool

	packets    []Entry
	packetList viewport.Model
	selected   int

	detailView viewport.Model

	focus FocusedPanel
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
		selected:   -1,
		focus:      NoFocus,
	}
}

func packetLines(packets []Entry, selected int) string {
	lines := make([]string, len(packets))
	for i, e := range packets {
		line := formatPacket(i+1, e, e.TimeStamp.Sub(packets[0].TimeStamp))
		if i == selected {
			line = selectedRowStyle.Render(line)
		}
		lines[i] = line
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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		l := m.layout()

		m.packetList.SetWidth(max(0, l.ContentWidth-borderSize))
		m.packetList.SetHeight(max(0, l.PacketsHeight-borderSize))

		m.detailView.SetWidth(max(0, l.ContentWidth-borderSize))
		m.detailView.SetHeight(max(0, l.DetailsHeight-borderSize))

		return m, nil

	case tea.MouseWheelMsg:
		if m.screenMode != Dashboard {
			return m, nil
		}

		switch m.focus {
		case PacketsFocus:
			m.packetList, cmd = m.packetList.Update(msg)
			return m, cmd
		case DetailsFocus:
			m.detailView, cmd = m.detailView.Update(msg)
			return m, cmd
		default:
			return m, nil
		}

	case tea.MouseClickMsg:
		if m.screenMode != Dashboard && msg.Button != tea.MouseLeft {
			return m, nil
		}

		row := msg.Y - m.packetPanelTop()
		if row < 0 || row >= m.packetList.Height() {
			return m, nil
		}

		index := m.packetList.YOffset() + row
		if index < len(m.packets) {
			m.selected = index
			m.packetList.SetContent(packetLines(m.packets, m.selected))
			m.detailView.SetContent(formatDetail(m.packets[m.selected]))
		}

		return m, nil

		// do here
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
			case "tab":
				if m.focus == PacketsFocus {
					m.focus = DetailsFocus
				} else {
					m.focus = PacketsFocus
				}
			default:
				m.packetList, cmd = m.packetList.Update(msg)
				return m, cmd
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
		m.packetList.SetContent(packetLines(m.packets, m.selected))
		if wasAtBottom {
			m.packetList.GotoBottom()
		}

		return m, handFrame(m.framesChan)
	}

	return m, nil
}
