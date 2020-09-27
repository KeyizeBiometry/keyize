package keyize

func AvgDynamics(d []*Dynamics) *Dynamics {
	propSet := newFloatSliceMapMan()

	for _, c := range d {
		for propName, p := range c.properties {
			propSet.Add(propName, p.Value)
		}
	}

	// Average values

	avgPropValues := propSet.Reduce()

	// Create the new Dynamics

	n := NewDynamics()

	for key, v := range avgPropValues {
		// There should be no error as the internally managed properties should be accurate and fully vetted
		prop, err := ParseDynamicsPropertyName(key)

		if err != nil {
			panic(err)
		}

		prop.Value = v

		n.AddProperty(prop)
	}

	return n
}
