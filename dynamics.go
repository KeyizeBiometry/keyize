// Package keyize provides abstractions to be used for the representation and computation of keystroke-related data.
// Keyize supports a wide array of keystroke recording and keystroke dynamics processing activities.
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
	// Consider properties from Dynamics d and a
	Both SharedPropertiesMethod = iota

	// Consider properties from only Dynamics d
	Left

	// Consider properties from only Dynamics a
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

// Properties returns the internally managed set of properties.
func (d *Dynamics) Properties() map[string]*DynamicsProperty {
	return d.properties
}

// AddProperty adds DynamicsProperty p to the internal map in Dynamics d.
func (d *Dynamics) AddProperty(p *DynamicsProperty) {
	d.properties[p.Name()] = p
}

// ProportionSharedProperties returns the proportion of properties shared between Dynamics d and a
// with respect to SharedPropertiesMethod method.
func (d *Dynamics) ProportionSharedProperties(a *Dynamics, method SharedPropertiesMethod) float64 {
	shared, total := d.SharedProperties(a, method)

	return float64(shared) / float64(total)
}

// SharedProperties returns the count of properties shared between Dynamics d and a
// as well as the count of total properties considered with respect to SharedPropertiesMethod method.
func (d *Dynamics) SharedProperties(a *Dynamics, method SharedPropertiesMethod) (shared int, total int) {
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
			if _, ok := targetB.properties[prop.Name()]; ok {
				shared++
			}

			total++
		}
	} else if method == Both {
		correspondingAPropsCount := 0

		for _, prop := range d.properties {
			if _, ok := a.properties[prop.Name()]; ok {
				shared++
				correspondingAPropsCount++
			}

			total++
		}

		// Add unaccounted A props to total.
		// It is known that these are not shared because we already found all shared props.
		total += len(a.properties) - correspondingAPropsCount
	}

	return shared, total
}

func (d *Dynamics) intermediateDist(a *Dynamics, squareDifferences bool) (dist float64) {
	td := 0.0

	for timingName, t1 := range d.properties {
		t2, ok := a.properties[timingName]

		if !ok {
			// Property not shared
			continue
		}

		// Scaling should be performed here

		if squareDifferences {
			td += math.Pow(math.Abs(t1.Value-t2.Value), 2)
		} else {
			td += math.Abs(t1.Value - t2.Value)
		}
	}

	return td
}

// ManhattanDist uses the Manhattan distance metric to find distance between Dynamics d and a.
func (d *Dynamics) ManhattanDist(a *Dynamics) (dist float64) {
	idist := d.intermediateDist(a, false)

	return idist
}

// EuclideanDist uses the Euclidean distance metric to find distance between Dynamics d and a.
func (d *Dynamics) EuclideanDist(a *Dynamics) (dist float64) {
	idist := d.intermediateDist(a, true)

	euclideanDist := math.Sqrt(idist)

	return euclideanDist
}
