# keyize
> Go Keyize keystroke dynamics library

[![PkgGoDev](https://pkg.go.dev/badge/github.com/KeyizeBiometry/keyize)](https://pkg.go.dev/github.com/KeyizeBiometry/keyize)

The keyize package provides types Recording and Dynamics.

Recording supports the manual addition of RecordingEvents, and can also parse recordings in the KeyizeV1 format. (eg. those created with [KeyizeBiometry/web-recorder](https://github.com/KeyizeBiometry/web-recorder))

Dynamics support a number of operations, and should cover most processing use-cases.

# Install

```sh
go get github.com/KeyizeBiometry/keyize
```

# Import Recording and Get Dynamics

```go
rec, err := keyize.ImportKeyizeV1("dH357uH582de1110ue1389dl2345ul2853dl3604ul4299do4458uo4729d 5583u 5815dw5942uw6196do6770uo7254dr7507ur7907dl8904ul9411dd10384ud10637")

// rec.Events will contain a list of parsed recording events

if err != nil {...}

// Convert to Dynamics

dyn := rec.Dynamics()
// Now an []DynamicsProperty has been extracted from the Recording and a Dynamics has been created
```

# Compare Dynamics

```go
// Proportion match (experimental)

manhattanDist := dyn1.ProportionMatch(dyn2)
// => 0.96

// Manhattan Distance

manhattanDist := dyn1.ManhattanDist(dyn2, nil)

// Euclidean Distance

euclideanDist := dyn1.EuclideanDist(dyn2, nil)

// Average property difference (scaled)

avgScaledDiff := dyn1.AvgScaledPropDiff(dyn2, nil)
```

# Note

This library is not yet complete. Some features are planned or being considered:
- [ ] Export Dynamics encoded as []byte using Protocol Buffers
- [ ] Export Recording encoded as KeyizeV1
- [ ] Refine and further test Dynamics ProportionMatch method
- [ ] Significantly improve test coverage
