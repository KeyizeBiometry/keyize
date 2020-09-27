package keyize

type floatSliceMapMan map[string][]float64

func newFloatSliceMapMan() floatSliceMapMan {
	return floatSliceMapMan{}
}

func (f floatSliceMapMan) Add(k string, v float64) {
	s, ok := f[k]

	if ok {
		s = append(s, v)
	}

	f[k] = []float64{v}
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
