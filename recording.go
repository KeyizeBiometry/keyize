package keyize

import (
	"errors"
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
	At      int
	Kind    RawEventKind
	Subject rune
}

// Text returns the (approximate) text typed in Recording r.
func (r *Recording) Text() string {
	var text []rune

	for _, e := range r.Events {
		if e.Kind == KeyDown {
			// Check if this is a special deletion character
			// We cannot do much with delete (\x7F) because we don't know cursor location when it is being used. Oh well.

			if e.Subject == '\b' {
				// Backspace. Delete last character.
				text = text[:len(text)-1]
			}

			// Otherwise append normally
			text = append(text, e.Subject)
		}
	}

	return string(text)
}

// Dynamics converts the raw data from Recording r to Dynamics d by extracting and averaging timings.
func (r *Recording) Dynamics() *Dynamics {
	// Create propTimings (prop -> timings slice)
	// First rune of key: a = DownDown, b = UpDown, c = Dwell

	propTimings := newFloatSliceMapMan()

	var lastDown rune
	var lastDownTime int

	var lastUp rune
	var lastUpTime int

	for i, e := range r.Events {
		switch e.Kind {
		case KeyDown:
			// DownDown / DD prop

			// => If last down was not zero value rune '\x00', push DD timing
			if lastDown != '\x00' {
				ddEventName := "a" + string(lastDown) + string(e.Subject)

				propTimings.Add(ddEventName, float64(e.At-lastDownTime))
			}

			// UpDown / UD prop

			// => If last up was not zero value rune '\x00', push UD timing
			if lastUp != '\x00' {
				udEventName := "b" + string(lastUp) + string(e.Subject)

				propTimings.Add(udEventName, float64(e.At-lastUpTime))
			}

			// Update lastDown

			lastDown = e.Subject
			lastDownTime = e.At
		case KeyUp:
			// Dwell prop

			// => Seek to last keydown for e.Subject if exists
			for b := i; b >= 0; b-- {
				relevantPreviousEvent := r.Events[b]

				if relevantPreviousEvent.Kind == KeyDown && relevantPreviousEvent.Subject == e.Subject {
					// Push timing and exit seeking

					dwellEventName := "c" + string(e.Subject)

					propTimings.Add(dwellEventName, float64(e.At-relevantPreviousEvent.At))

					break
				}
			}

			// Update lastUp

			lastUp = e.Subject
			lastUpTime = e.At
		}
	}

	// Average propTimings

	avgPropTimings := propTimings.Reduce()

	// Create the Dynamics and convert avgPropTimings timings to DynamicsProperties

	d := NewDynamics()

	for propKey, avgTime := range avgPropTimings {
		// It may be assumed that first rune is valid because it was created above

		propTypeRune, _ := utf8.DecodeRuneInString(propKey)

		switch propTypeRune {
		case 'a':
			// DownDown

			fallthrough
		case 'b':
			// UpDown

			rune1, _ := utf8.DecodeRuneInString(propKey[1:])
			rune2, _ := utf8.DecodeRuneInString(propKey[2:])

			propKind := UpDown

			if propTypeRune == 'a' {
				propKind = DownDown
			}

			d.AddProperty(&DynamicsProperty{
				Kind:  propKind,
				KeyA:  rune1,
				KeyB:  rune2,
				Value: avgTime,
			})
		case 'c':
			// Dwell

			rune1, _ := utf8.DecodeRuneInString(propKey[1:])

			d.AddProperty(&DynamicsProperty{
				Kind:  Dwell,
				KeyA:  rune1,
				KeyB:  '\x00',
				Value: avgTime,
			})
		}
	}

	return d
}

// ImportKeyizeV1 imports a keystroke recording of the Keyize V1 format.
// These recordings may be generated from users by a corresponding recording library.
func ImportKeyizeV1(d string) (*Recording, error) {
	rec := &Recording{
		Events: []*RecordingEvent{},
	}

	matches := v1regex.FindAllStringSubmatch(d, -1)

	var lastAt int64 = -100

	for _, m := range matches {
		kindRune, _ := utf8.DecodeRuneInString(m[1])
		subjectRune, _ := utf8.DecodeRuneInString(m[2])

		if kindRune == utf8.RuneError || subjectRune == utf8.RuneError {
			return nil, errors.New("invalid subject or kind rune")
		}

		eventKind, ok := runeEventKindMap[kindRune]

		if !ok {
			return nil, errors.New("invalid event kind " + string(kindRune))
		}

		at, err := strconv.ParseInt(m[3], 10, 64)

		if err != nil {
			return nil, err
		}

		if at < lastAt {
			return nil, errors.New("invalid at value " + strconv.Itoa(int(at)) + " is less than previous")
		}

		if at < 0 {
			return nil, errors.New("invalid at value " + strconv.Itoa(int(at)) + " is less than 0")
		}

		rec.Events = append(rec.Events, &RecordingEvent{
			Kind:    eventKind,
			At:      int(at),
			Subject: subjectRune,
		})
	}

	return rec, nil
}
