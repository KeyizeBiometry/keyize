package main

import (
	"encoding/csv"
	"fmt"
	"github.com/KeyizeBiometry/keyize"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type session struct {
	subject string
	d       *keyize.Dynamics
}

var firstSegKindMap = map[string]keyize.DynamicsPropertyKind{
	"H":  keyize.Dwell,
	"DD": keyize.DownDown,
	"UD": keyize.UpDown,
}

var convRuneMap = map[string]rune{
	"Shift":  '\x0F',
	"period": '.',
	"five":   '5',
	"Return": '\n',
}

func convRune(s string) rune {
	if len(s) == 1 {
		r, _ := utf8.DecodeRuneInString(s)
		return r
	} else {
		return convRuneMap[s]
	}
}

func avgFloatSlice(a []float64) float64 {
	t := 0.0

	for _, v := range a {
		t += v
	}

	return t / float64(len(a))
}

var sessions []*session

func main() {
	// Open file

	f, err := os.Open("DSL-StrongPasswordData.csv")

	if err != nil {
		if err == os.ErrNotExist {
			fmt.Println("Tip: The DSL-StrongPasswordData.csv dataset may be downloaded from https://www.cs.cmu.edu/~keystroke/")
		}

		panic(err)
	}

	// Parse data to memory (sessions)

	fmt.Println("Loading data...")

	r := csv.NewReader(f)

	var fieldOrder []string

	for i := 0; ; i++ {
		record, err := r.Read()

		if err != nil {
			break
		}

		if i == 0 {
			// Get field order from first record

			fieldOrder = record
			continue
		}

		s := &session{
			subject: "",
			d:       keyize.NewDynamics(),
		}

		for i, v := range record {
			fieldName := fieldOrder[i]

			if fieldName == "subject" {
				s.subject = v
			} else if fieldName != "sessionIndex" && fieldName != "rep" {
				fieldNameSegs := strings.Split(fieldName, ".")

				kind := firstSegKindMap[fieldNameSegs[0]]

				keyA := convRune(fieldNameSegs[1])
				var keyB rune

				if len(fieldNameSegs) >= 3 {
					keyB = convRune(fieldNameSegs[2])
				}

				secondsValue, err := strconv.ParseFloat(v, 64)

				if err != nil {
					panic(err)
				}

				s.d.AddProperty(&keyize.DynamicsProperty{
					Kind:  kind,
					KeyA:  keyA,
					KeyB:  keyB,
					Value: secondsValue * 1000,
				})
			}
		}

		sessions = append(sessions, s)
	}

	// Now summarize data into users using averages

	fmt.Println("Done loading")

	performClosedSetID()
}

func performClosedSetID() {
	fmt.Println("== Closed Set Identification Processing ==")

	correct := 0
	incorrect := 0

	for i, a := range sessions {
		fmt.Println(i, "/", len(sessions))

		var bestDist float64
		var bestMatch *session = nil

		for _, b := range sessions {
			// Ignore if sessions are same session
			if a == b {
				continue
			}

			dist, err := a.d.ManhattanDist(b.d, nil)

			if err != nil {
				panic(err)
			}

			if bestMatch == nil || dist < bestDist {
				bestMatch = b
				bestDist = dist
			}
		}

		if a.subject == bestMatch.subject {
			correct++
		} else {
			incorrect++
		}
	}

	fmt.Printf("Accuracy: %.2f", 100*float64(correct)/float64(correct+incorrect))
}
