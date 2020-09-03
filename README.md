# keyize
> Go Keyize keystroke dynamics library

Planned additions:
- Possibly add DynamicsGroup type for closed set identification? Consider how this might be useful / not useful.
- Add Compare function to Dynamics which returns a percent match, rather than raw distance, and a confidence from the raw distance method?
- Add intermediaryDist function to Dynamics which returns total property distance, with a parameter which can be used to square the differences while summing. This function will be used as the intermediary for both EuclideanDist and ManhattanDist, and can perform scaling based on property type. Consider changing the properties keys to not just float64, so they can be identified for scaling faster.
