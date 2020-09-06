package keyize

import (
	"testing"
)

func TestDynamics_GetProportionSharedProperties(t *testing.T) {
	dynA := NewDynamics()
	dynB := NewDynamics()

	dynA.AddProperty(&DynamicsProperty{
		Kind:  DownDown,
		KeyA:  'a',
		KeyB:  'b',
		Value: 4.55,
	})

	dynA.AddProperty(&DynamicsProperty{
		Kind:  UpDown,
		KeyA:  'b',
		KeyB:  'd',
		Value: 8.88,
	})

	dynB.AddProperty(&DynamicsProperty{
		Kind:  DownDown,
		KeyA:  'a',
		KeyB:  'b',
		Value: 6.77,
	})

	dynB.AddProperty(&DynamicsProperty{
		Kind:  UpDown,
		KeyA:  'b',
		KeyB:  'c',
		Value: 3.33,
	})

	dynB.AddProperty(&DynamicsProperty{
		Kind:  UpDown,
		KeyA:  'c',
		KeyB:  'd',
		Value: 7.33,
	})

	// getPropAssertMultidirectional executes the requested GetProportionSharedProperties as well as the appropriate inverse to ensure that expected values are produced with both a and b as the argument Dynamics.
	getPropAssertMultidirectional := func(a *Dynamics, b *Dynamics, method SharedPropertiesMethod) float64 {
		var methodB SharedPropertiesMethod

		// Select appropriate inverse method

		switch method {
		case Both:
			methodB = Both
		case Left:
			methodB = Right
		case Right:
			methodB = Left
		}

		// Execute both

		resA := a.GetProportionSharedProperties(dynB, method)
		resB := b.GetProportionSharedProperties(dynA, methodB)

		// Ensure requested is equal to inverse

		if resA != resB {
			t.Fatalf("could not assert multidirectional %d %d: %f != %f", method, methodB, resA, resB)
		}

		t.Log("Passed multidirectional", resA, resB)

		return resA
	}

	if v := getPropAssertMultidirectional(dynA, dynB, Both); v != 0.25 {
		t.Errorf("bad Both %f", v)
	}

	if v := getPropAssertMultidirectional(dynA, dynB, Right); v < 0.33 || v > 0.34 {
		t.Errorf("bad Right %f", v)
	}

	if v := getPropAssertMultidirectional(dynA, dynB, Left); v != 0.5 {
		t.Errorf("bad Left %f", v)
	}
}
