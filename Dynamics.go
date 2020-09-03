package keyize

import "math"

// Dynamics represents a group of dynamics properties
type Dynamics struct {
	// timings is an internally managed set of dynamics properties
	// It is private to prevent mismanagement, as the keys should be
	// formatted correctly.
	timings map[string]float64
}

func (d *Dynamics) ManhattanDist(a *Dynamics) (match float64, confidence float64) {
	td := 0.0
	timingsNotInCommon := 0
	timingsInCommon := 0

	for timingName, t1 := range d.timings {
		t2, ok := a.timings[timingName]

		if !ok {
			timingsNotInCommon++
			continue
		}

		timingsInCommon++
		td += math.Abs(t1 - t2)
	}

	totalTimings := timingsNotInCommon + timingsInCommon

	return td, float64(timingsInCommon) / (float64(totalTimings))
}

func (d *Dynamics) EuclideanDist(a *Dynamics) (match float64, confidence float64) {
	td := 0.0
	timingsNotInCommon := 0
	timingsInCommon := 0

	for timingName, t1 := range d.timings {
		t2, ok := a.timings[timingName]

		if !ok {
			timingsNotInCommon++
			continue
		}

		timingsInCommon++
		td += math.Pow(math.Abs(t1-t2), 2)
	}

	euclideanDist := math.Sqrt(td)

	totalTimings := timingsNotInCommon + timingsInCommon

	return euclideanDist, float64(timingsInCommon) / (float64(totalTimings))
}
