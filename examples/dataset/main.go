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

type subject struct {
	name string
	sessions   []*keyize.Dynamics
	avgSession *keyize.Dynamics
}

var subjects = map[string]*subject{}

func main() {
	// Open file

	f, err := os.Open("DSL-StrongPasswordData.csv")

	if err != nil {
		if err == os.ErrNotExist {
			fmt.Println("Tip: The DSL-StrongPasswordData.csv dataset may be downloaded from https://www.cs.cmu.edu/~keystroke/")
		}

		panic(err)
	}

	defer f.Close()

	// Parse data to memory (subjects)

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

		if len(subjects) >= 10 || i < 10 {
			continue
		}

		s := keyize.NewDynamics()
		var curSubj string

		for i2, v := range record {
			fieldName := fieldOrder[i2]

			if fieldName == "subject" {
				curSubj = v
			} else if fieldName != "sessionIndex" && fieldName != "rep" {
				fieldNameSegs := strings.Split(fieldName, ".")

				if fieldNameSegs[0] == "H" {
					fieldNameSegs[0] = "D"
				}

				finalFieldName := strings.Join(fieldNameSegs, ".")

				secondsValue, err := strconv.ParseFloat(v, 64)

				if err != nil {
					panic(err)
				}

				s.AddPropertyByName(finalFieldName, secondsValue * 1000)
			}
		}

		relevSubj, ok := subjects[curSubj]

		if !ok {
			subjects[curSubj] = &subject{
				name: curSubj,
				sessions:   []*keyize.Dynamics{s},
				avgSession: nil,
			}
		} else {
			relevSubj.sessions = append(relevSubj.sessions, s)
		}
	}

	// Now summarize data for subject averages

	for _, u := range subjects {
		u.avgSession = keyize.AvgDynamics(u.sessions)
	}

	// Done. Perform tests.

	fmt.Println("Done loading", len(subjects), "subjects.")

	performClosedSetID()
}

func performClosedSetID() {
	fmt.Println("== Closed Set Identification Processing ==")

	correct := 0
	incorrect := 0

	for topSubjName, topSubj := range subjects {
		for _, cdyn := range topSubj.sessions {
			var bestMatchSubjName string
			var bestMatchDist float64 = -2

			for csubjName, csubj := range subjects {
				dist := cdyn.ManhattanDist(csubj.avgSession, nil)

				if dist < bestMatchDist || bestMatchDist == -2 {
					bestMatchSubjName = csubjName
					bestMatchDist = dist
				}
			}

			if bestMatchSubjName == topSubjName {
				correct++
			} else {
				//fmt.Println(bestMatchSubjName, topSubjName, bestMatchDist)
				incorrect++
			}
		}
	}

	fmt.Printf("Accuracy: %.2f (%d / %d)", 100*float64(correct)/float64(correct+incorrect), correct, correct+incorrect)
}
