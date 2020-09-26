# keyize
> Go Keyize keystroke dynamics library

Planned additions:
- Possibly add DynamicsGroup type for closed set identification? Consider how this might be useful / not useful.
- Add Compare function to Dynamics which returns a percent match, rather than raw distance, and confidence from ProportionSharedProperties? Instead of dist using Dist methods, use to-be-implemented AvgScaledPropDist.
- Add numericMapMan for common reduction / map append implementation
- Add AvgScaledPropDist for Dynamics. Test using defaultDynamicsPropertyKindScaleMap.
