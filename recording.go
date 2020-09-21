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

func (r *Recording) GetText() string {
	var text []rune

	for _, e := range r.Events {
		if e.Kind == KeyDown {
			text = append(text, e.Subject)
		}
	}

	return string(text)
}

// pushEvent adds an int value to an []int value of key eventName in map m.
// If the []int does not exist at key eventName, pushEvent will initialize it.
func pushEvent(m map[string][]int, eventName string, value int) {
	_, ok := m[eventName]

	if !ok {
		m[eventName] = []int{value}
	} else {
		m[eventName] = append(m[eventName], value)
	}
}

// ToDynamics converts the raw data from Recording r to Dynamics d by extracting and averaging timings.
func (r *Recording) ToDynamics() *Dynamics {
	// Create propTimings (prop -> timings slice)
	// First rune of key: a = DownDown, b = UpDown, c = Dwell

	propTimings := map[string][]int{}

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

				pushEvent(propTimings, ddEventName, e.At-lastDownTime)
			}

			// UpDown / UD prop

			// => If last up was not zero value rune '\x00', push UD timing
			if lastUp != '\x00' {
				udEventName := "b" + string(lastUp) + string(e.Subject)

				pushEvent(propTimings, udEventName, e.At-lastUpTime)
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

					pushEvent(propTimings, dwellEventName, e.At-relevantPreviousEvent.At)

					break
				}
			}

			// Update lastUp

			lastUp = e.Subject
			lastUpTime = e.At
		}
	}

	// Average propTimings

	avgPropTimings := map[string]int{}

	for prop, timings := range propTimings {
		total := 0

		for _, t := range timings {
			total += t
		}

		avgPropTimings[prop] = total / len(timings)
	}

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
				Value: float64(avgTime),
			})
		case 'c':
			// Dwell

			rune1, _ := utf8.DecodeRuneInString(propKey[1:])

			d.AddProperty(&DynamicsProperty{
				Kind:  Dwell,
				KeyA:  rune1,
				KeyB:  '\x00',
				Value: float64(avgTime),
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

		rec.Events = append(rec.Events, &RecordingEvent{
			Kind:    eventKind,
			At:      int(at),
			Subject: subjectRune,
		})
	}

	return rec, nil
}
