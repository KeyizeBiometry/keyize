package keyize

import (
	"testing"
)

func TestAvgDynamics(t *testing.T) {
	d1 := NewDynamics()
	d2 := NewDynamics()

	d1.AddPropertyByName("D.A", 2)
	d2.AddPropertyByName("D.A", 4)

	d1.AddPropertyByName("DD.B.C", 2)
	d2.AddPropertyByName("DD.B.C", 4)

	d1.AddPropertyByName("UD.D.e", 4.5)
	d2.AddPropertyByName("UD.D.e", 4.7)

	d1.AddPropertyByName("UD.F.G", 4.2)

	avg := AvgDynamics([]*Dynamics{d1, d2})

	avgProps := avg.Properties()

	if avgProps["D.A"].Value != 3 || avgProps["DD.B.C"].Value != 3 || avgProps["UD.D.e"].Value != 4.6 || avgProps["UD.F.G"].Value != 4.2 {
		for _, p := range avgProps {
			t.Log(p.Name(), p.Value)
		}

		t.Fatal("Incorrect average property values")
	}
}

func TestDynamics_AddPropertyByName(t *testing.T) {
	d := NewDynamics()

	d.AddPropertyByName("D.H", 5)
	d.AddPropertyByName("DD.j.W", 4)
	d.AddPropertyByName("UD.o.E", 3)

	props := d.Properties()

	if props["D.H"].Value != 5 || props["DD.j.W"].Value != 4 || props["UD.o.E"].Value != 3 {
		t.Fatal("Not all properties exist / some may have incorrect values")
	} else {
		t.Log("Confirmed correct properties were added")
	}

	if props["D.H"].KeyA != 'H' || props["DD.j.W"].KeyA != 'j' || props["DD.j.W"].KeyB != 'W' || props["UD.o.E"].KeyA != 'o' || props["UD.o.E"].KeyB != 'E' {
		t.Fatal("Properties have incorrect KeyA / KeyB values")
	} else {
		t.Log("Confirmed correct KeyA / KeyB were determined")
	}
}

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

	// getPropAssertMultidirectional executes the requested ProportionSharedProperties as well as the appropriate inverse to ensure that expected values are produced with both a and b as the argument Dynamics.
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

		resA := a.ProportionSharedProperties(dynB, method)
		resB := b.ProportionSharedProperties(dynA, methodB)

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
