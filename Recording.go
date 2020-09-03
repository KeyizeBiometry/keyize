package keyize

import "time"

type RawEventKind int

const (
	KeyDown RawEventKind = iota
	KeyUp
)

// Recording represents a user's raw typing recording
type Recording struct {
	Events []RecordingEvent
}

// RecordingEvent is a specific event which took place during a recording
type RecordingEvent struct {
	At      time.Time
	Kind    RawEventKind
	Subject rune
}

// TODO: add import method for Recording
// func (r *Recording) ImportRaw(raw string) error

// Or: add new recording from raw function
// func ImportRaw(raw string) (*Recording, error)
