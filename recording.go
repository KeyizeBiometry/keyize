package keyize

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"unicode/utf8"
)

type RawEventKind int

const (
	KeyDown RawEventKind = iota
	KeyUp
)

var runeEventKindMap = map[rune]RawEventKind{
	'u': KeyUp,
	'd': KeyDown,
}

var v1regex = regexp.MustCompile("([a-z])(.|\n)([0-9]+)")

// Recording represents a user's raw typing recording
type Recording struct {
	Events []*RecordingEvent
}

// RecordingEvent is a specific event which took place during a recording
type RecordingEvent struct {
	At      uint64
	Kind    RawEventKind
	Subject rune
}

// ImportKeyizeV1 imports a keystroke recording of the Keyize V1 format.
// These recordings may be generated from users by a corresponding recording library.
func ImportKeyizeV1(d string) (*Recording, error) {
	rec := &Recording{
		Events: []*RecordingEvent{},
	}

	matches := v1regex.FindAllStringSubmatch(d, -1)

	for _, m := range matches {
		fmt.Println(m)

		kindRune, _ := utf8.DecodeRuneInString(m[1])
		subjectRune, _ := utf8.DecodeRuneInString(m[2])

		if kindRune == utf8.RuneError || subjectRune == utf8.RuneError {
			return nil, errors.New("invalid subject or kind rune")
		}

		eventKind, ok := runeEventKindMap[kindRune]

		if !ok {
			return nil, errors.New("invalid event kind " + string(kindRune))
		}

		at, err := strconv.ParseUint(m[3], 10, 64)

		if err != nil {
			return nil, err
		}

		rec.Events = append(rec.Events, &RecordingEvent{
			Kind:    eventKind,
			At:      at,
			Subject: subjectRune,
		})
	}

	return rec, nil
}

// TODO: add import method for Recording
// func (r *Recording) ImportRaw(raw string) error

// Or: add new recording from raw function
// func ImportRaw(raw string) (*Recording, error)
