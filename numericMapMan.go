package keyize

// Note for the future: once vanilla Go supports generics, floatSliceMapMan should be converted
// to numericSliceMapMan. Then in recording.Dynamics(), ints should be used for most calculations until the final
// creation of a Dynamics.
//
// This will allow the avoidance of unnecessary float math operations.

type floatSliceMapMan map[string][]float64

func newFloatSliceMapMan() floatSliceMapMan {
	return floatSliceMapMan{}
}

func (f floatSliceMapMan) Add(k string, v float64) {
	s, ok := f[k]

	if ok {
		f[k] = append(s, v)
	} else {
		f[k] = []float64{v}
	}
}

func (f floatSliceMapMan) Reduce() map[string]float64 {
	reduced := map[string]float64{}

	for k, vs := range f {
		t := 0.0

		for _, v := range vs {
			t += v
		}

		reduced[k] = t / float64(len(vs))
	}

	return reduced
}
