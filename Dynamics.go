package keyize

import "math"

type DynamicsPropertyKind int

const (
	Dwell DynamicsPropertyKind = iota
	DownDown
	UpDown
)

// DynamicsProperty is a keystroke dynamics property, often synthesized from
// a Recording.
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
	// properties is an internally managed set of dynamics properties
	// It is private to prevent the introduction of malformatted keys.
	properties map[string]*DynamicsProperty
}

// AddProperty adds DynamicsProperty p to the internal map in Dynamics d.
func (d *Dynamics) AddProperty(p *DynamicsProperty) {
	d.properties[p.Name()] = p
}

func (d *Dynamics) intermediateDist(a *Dynamics, squareDifferences bool) (dist float64, confidence float64) {
	td := 0.0
	propsNotInCommonCount := 0
	propsInCommonCount := 0

	for timingName, t1 := range d.properties {
		t2, ok := a.properties[timingName]

		if !ok {
			propsNotInCommonCount++
			continue
		}

		// Scaling should be performed here

		propsInCommonCount++

		if squareDifferences {
			td += math.Pow(math.Abs(t1.Value-t2.Value), 2)
		} else {
			td += math.Abs(t1.Value - t2.Value)
		}
	}

	totalTimings := propsNotInCommonCount + propsInCommonCount

	return td, float64(propsInCommonCount) / (float64(totalTimings))
}

// ManhattanDist uses the Manhattan distance metric to find distance between Dynamics d and a.
// It also provides a confidence level, which represents the proportion of shared, comparable properties.
func (d *Dynamics) ManhattanDist(a *Dynamics) (dist float64, confidence float64) {
	idist, confidence := d.intermediateDist(a, false)

	return idist, confidence
}

// EuclideanDist uses the Euclidean distance metric to find distance between Dynamics d and a.
// It also provides a confidence level, which represents the proportion of shared, comparable properties.
func (d *Dynamics) EuclideanDist(a *Dynamics) (dist float64, confidence float64) {
	idist, confidence := d.intermediateDist(a, false)

	euclideanDist := math.Sqrt(idist)

	return euclideanDist, confidence
}
