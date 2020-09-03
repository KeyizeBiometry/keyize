# keyize
> Go Keyize keystroke dynamics library

Planned additions:
- Possibly add DynamicsGroup type for closed set identification? Consider how this might be useful / not useful.
- Add Compare function to Dynamics which returns a percent match, rather than raw distance, and a confidence from the raw distance method?
- Use intermediateDist to perform scaling based on property type. Potentially update values in Dynamics timings to something which could make the property type identification process faster than by string key processing.
