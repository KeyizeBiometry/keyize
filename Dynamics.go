package keyize

import (
	"math"
)

type DynamicsPropertyKind int

const (
	Dwell DynamicsPropertyKind = iota
	DownDown
	UpDown
)

type SharedPropertiesMethod int

const (
	Both SharedPropertiesMethod = iota
	Left
	Right
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

func NewDynamics() *Dynamics {
	return &Dynamics{
		properties: map[string]*DynamicsProperty{},
	}
}

// AddProperty adds DynamicsProperty p to the internal map in Dynamics d.
func (d *Dynamics) AddProperty(p *DynamicsProperty) {
	d.properties[p.Name()] = p
}

func (d *Dynamics) GetProportionSharedProperties(a *Dynamics, method SharedPropertiesMethod) float64 {
	shared := 0
	total := 0

	if method == Right || method == Left {
		var targetA *Dynamics
		var targetB *Dynamics

		if method == Right {
			targetA = a
			targetB = d
		} else {
			targetA = d
			targetB = a
		}

		for _, prop := range targetA.properties {
			_, ok := targetB.properties[prop.Name()]

			if ok {
				shared++
			}

			total++
		}
	} else if method == Both {
		discoveredAProps := []*DynamicsProperty{}

		for _, prop := range d.properties {
			p, ok := a.properties[prop.Name()]

			if ok {
				shared++
				discoveredAProps = append(discoveredAProps, p)
			}

			total++
		}

		for _, prop := range a.properties {
			// Check if prop is already discovered
			alreadyDiscovered := false

			for _, p := range discoveredAProps {
				if prop == p {
					alreadyDiscovered = true
				}
			}

			if !alreadyDiscovered {
				// If not already discovered, this property is also not shared.
				total++
			}
		}
	}

	return float64(shared) / float64(total)
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

	totalProps := propsNotInCommonCount + propsInCommonCount

	return td, float64(propsInCommonCount) / (float64(totalProps))
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
	idist, confidence := d.intermediateDist(a, true)

	euclideanDist := math.Sqrt(idist)

	return euclideanDist, confidence
}
