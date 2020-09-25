// Package keyize provides abstractions to be used for the representation and computation of keystroke-related data.
// Keyize supports a wide array of keystroke recording and keystroke dynamics processing activities.
package keyize

import (
	"errors"
	"math"
	"strconv"
)

type DynamicsPropertyKind int

const (
	Dwell DynamicsPropertyKind = iota
	DownDown
	UpDown
)

// defaultDynamicsPropertyKindScaleMap is the optimized kind scale map discovered during research
var defaultDynamicsPropertyKindScaleMap = DynamicsPropertyKindScaleMap{
	Dwell:    1,
	DownDown: 1 / 19.4,
	UpDown:   1 / 14,
}

// DynamicsPropertyKindScaleMap is a map of DynamicsPropertyKind to scaling values.
// It is used when calculating distance between multiple Dynamics.
//
// If a key is missing, it may be replaced with the default, optimized value if present.
// Using the optimized default scale values by passing nil to distance methods for the propertyKindScaleMap argument is recommended.
type DynamicsPropertyKindScaleMap map[DynamicsPropertyKind]float64

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

// RemoveProperty removes DynamicsProperty p from the internal map in Dynamics d.
func (d *Dynamics) RemoveProperty(p *DynamicsProperty) {
	delete(d.properties, p.Name())
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

func (d *Dynamics) intermediateDist(a *Dynamics, squareDifferences bool, propertyKindScaleMap DynamicsPropertyKindScaleMap) (dist float64, err error) {
	td := 0.0

	for timingName, t1 := range d.properties {
		t2, ok := a.properties[timingName]

		if !ok {
			// Property not shared
			continue
		}

		// Perform scaling

		scaleProp, ok := propertyKindScaleMap[t1.Kind]

		if !ok {
			scaleProp, ok = defaultDynamicsPropertyKindScaleMap[t1.Kind]

			if !ok {
				// Default map and provided map do not include scale for kind

				return 900, errors.New("could not locate a scale value for kind " + strconv.Itoa(int(t1.Kind)))
			}
		}

		scaledT1Value := t1.Value * scaleProp
		scaledT2Value := t2.Value * scaleProp

		// Calculate total using scaled values

		if squareDifferences {
			td += math.Pow(math.Abs(scaledT1Value-scaledT2Value), 2)
		} else {
			td += math.Abs(scaledT1Value - scaledT2Value)
		}
	}

	return td, nil
}

// ManhattanDist uses the Manhattan distance metric to find distance between Dynamics d and a.
//
// It uses propertyKindScaleMap for scaling. nil may be passed for propertyKindScaleMap and the optimized defaults will be used.
func (d *Dynamics) ManhattanDist(a *Dynamics, propertyKindScaleMap DynamicsPropertyKindScaleMap) (dist float64, err error) {
	idist, err := d.intermediateDist(a, false, propertyKindScaleMap)

	return idist, err
}

// EuclideanDist uses the Euclidean distance metric to find distance between Dynamics d and a.
//
// It uses propertyKindScaleMap for scaling. nil may be passed for propertyKindScaleMap and the optimized defaults will be used.
func (d *Dynamics) EuclideanDist(a *Dynamics, propertyKindScaleMap DynamicsPropertyKindScaleMap) (dist float64, err error) {
	idist, err := d.intermediateDist(a, true, propertyKindScaleMap)

	euclideanDist := math.Sqrt(idist)

	return euclideanDist, err
}
