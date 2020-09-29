package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/KeyizeBiometry/keyize"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf8"
)

var maxGroupSize *int = flag.Int("maxClosedSetGroupSize", 10, "Max group size for closed set identification")

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
	name       string
	sessions   []*keyize.Dynamics
	avgSession *keyize.Dynamics
}

var subjects = []*subject{}

func main() {
	flag.Parse()

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

				s.AddPropertyByName(finalFieldName, secondsValue*1000)
			}
		}

		sameSubjectIdx := -1

		for idx, s := range subjects {
			if s.name == curSubj {
				sameSubjectIdx = idx
			}
		}

		if sameSubjectIdx == -1 {
			subjects = append(subjects, &subject{
				name:       curSubj,
				sessions:   []*keyize.Dynamics{s},
				avgSession: nil,
			})
		} else {
			subjects[sameSubjectIdx].sessions = append(subjects[sameSubjectIdx].sessions, s)
		}
	}

	// Now summarize data for subject averages

	for _, u := range subjects {
		u.avgSession = keyize.AvgDynamics(u.sessions)
	}

	// Done. Perform tests.

	fmt.Println("Done loading", len(subjects), "subjects.")

	StartTime("ClosedSetID")

	performClosedSetID()

	EndTime()
}

func groupClosedSetIDStats(subjects []*subject) float64 {
	var correct int64
	var incorrect int64

	wg := &sync.WaitGroup{}

	for _, topSubj := range subjects {
		wg.Add(1)

		go func() {
			for _, cdyn := range topSubj.sessions {
				var bestMatchSubj *subject
				var bestMatchDist float64 = -2

				for _, csubj := range subjects {
					dist := cdyn.ManhattanDist(csubj.avgSession, nil)

					if dist < bestMatchDist || bestMatchDist == -2 {
						bestMatchSubj = csubj
						bestMatchDist = dist
					}
				}

				if bestMatchSubj == topSubj {
					atomic.AddInt64(&correct, 1)
				} else {
					atomic.AddInt64(&incorrect, 1)
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()

	return 100 * float64(correct) / float64(correct+incorrect)
}

func performClosedSetID() {
	fmt.Println("== Closed Set Identification Processing ==")

	// For each group size

	for groupSize := 2; groupSize <= *maxGroupSize; groupSize++ {
		// For each possible sequential group in subjects

		var results []float64

		for startIdx := 0; startIdx < len(subjects)-groupSize; startIdx++ {
			useGroup := subjects[startIdx : startIdx+groupSize]

			results = append(results, groupClosedSetIDStats(useGroup))
		}

		// Average results

		t := 0.0

		for _, r := range results {
			t += r
		}

		avgAccuracy := t / float64(len(results))

		fmt.Printf("GS %d   %d   %.2f\n", groupSize, len(results), avgAccuracy)
	}
}
