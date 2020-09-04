package keyize

import "math"

type DynamicsPropertyKind int

const (
	Dwell DynamicsPropertyKind = iota
	DownDown
	UpDown
)

type DynamicsProperty struct {
	Kind  DynamicsPropertyKind
	KeyA  rune
	KeyB  rune
	Value float64
}

// Name provides the name for the property (this can be used as a key by a Dynamics)
func (d *DynamicsProperty) Name() string {
	switch d.Kind {
	case Dwell:
		return "D." + string(d.KeyA)
	case DownDown:
		return "DD." + string(d.KeyA) + "." + string(d.KeyB)
	case UpDown:
		return "UD." + string(d.KeyA) + "." + string(d.KeyB)
	default:
		panic("cannot create a name for a DynamicsProperty of unknown kind")
	}
}

// Dynamics represents a group of dynamics properties
type Dynamics struct {
	// timings is an internally managed set of dynamics properties
	// It is private to prevent the introduction of malformatted keys.
	timings map[string]*DynamicsProperty
}

func (d *Dynamics) PushProperty(p *DynamicsProperty) {
	d.timings[p.Name()] = p
}

func (d *Dynamics) intermediateDist(a *Dynamics, squareDifferences bool) (dist float64, confidence float64) {
	td := 0.0
	timingsNotInCommon := 0
	timingsInCommon := 0

	for timingName, t1 := range d.timings {
		t2, ok := a.timings[timingName]

		if !ok {
			timingsNotInCommon++
			continue
		}

		// Scaling should be performed here

		timingsInCommon++

		if squareDifferences {
			td += math.Pow(math.Abs(t1.Value-t2.Value), 2)
		} else {
			td += math.Abs(t1.Value - t2.Value)
		}
	}

	totalTimings := timingsNotInCommon + timingsInCommon

	return td, float64(timingsInCommon) / (float64(totalTimings))
}

func (d *Dynamics) ManhattanDist(a *Dynamics) (match float64, confidence float64) {
	idist, confidence := d.intermediateDist(a, false)

	return idist, confidence
}

func (d *Dynamics) EuclideanDist(a *Dynamics) (match float64, confidence float64) {
	idist, confidence := d.intermediateDist(a, false)

	euclideanDist := math.Sqrt(idist)

	return euclideanDist, confidence
}
