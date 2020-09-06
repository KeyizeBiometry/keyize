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
		Kind:  DownDown,
		KeyA:  'h',
		KeyB:  'j',
		Value: 8.88,
	})

	dynB.AddProperty(&DynamicsProperty{
		Kind:  DownDown,
		KeyA:  'a',
		KeyB:  'b',
		Value: 6.77,
	})

	dynB.AddProperty(&DynamicsProperty{
		Kind:  DownDown,
		KeyA:  'k',
		KeyB:  'l',
		Value: 3.33,
	})

	getPropAssertMultidirectional := func(a *Dynamics, b *Dynamics, method SharedPropertiesMethod) float64 {
		var methodB SharedPropertiesMethod

		switch method {
		case Both:
			methodB = Both
		case Left:
			methodB = Right
		case Right:
			methodB = Left
		}

		resA := a.GetProportionSharedProperties(dynB, method)
		resB := b.GetProportionSharedProperties(dynA, methodB)

		if resA != resB {
			t.Fatalf("could not assert multidirectional %d %d: %f != %f", method, methodB, resA, resB)
		}

		t.Log("Passed multidirectional", resA, resB)

		return resA
	}

	if v := getPropAssertMultidirectional(dynA, dynB, Both); v < 0.33 || v > 0.34 {
		t.Errorf("bad Both %f", v)
	}

	if v := getPropAssertMultidirectional(dynA, dynB, Right); v != 0.5 {
		t.Errorf("bad Right %f", v)
	}

	if v := getPropAssertMultidirectional(dynA, dynB, Left); v != 0.5 {
		t.Errorf("bad Left %f", v)
	}
}
