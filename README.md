# mako_time_converter (go)
![Unittest status badge](https://github.com/hochfrequenz/mako_time_converter/workflows/Unittests/badge.svg)
![Coverage status badge](https://github.com/hochfrequenz/mako_time_converter/workflows/coverage/badge.svg)
![Linter status badge](https://github.com/hochfrequenz/mako_time_converter/workflows/golangci-lint/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/hochfrequenz/mako_time_converter.svg)](https://pkg.go.dev/github.com/hochfrequenz/mako_time_converter)

This is a Go module to convert between German "Gastag" and "Stromtag", between inclusive and exclusive end dates and combinations of all them.
This is relevant for German Marktkommunikation ("MaKo").

This is a Golang port of the [`mako_datetime_converter` for .NET](https://github.com/Hochfrequenz/mako_datetime_converter).

## Rationale

The German Marktkommunikation ("MaKo") defines some rules regarding date times:

- you shall communicate end dates as exclusive (which is [generally a good idea](https://hf-kklein.github.io/exclusive_end_dates.github.io/))
- you shall use UTC date times with a specified UTC offset (which is a good idea)
- and you shall always use UTC-offset 0 (which makes things unnecessary complicated)
- in electricity all days (and contracts) start and end at midnight of German local time
- but in gas all days (and contracts) start and end at 6am German local time ("Gas-Tag")

Now imagine there is an interface between two systems:

- one of your systems obeys all of the above rules
- but another one works differently (e.g. models end dates inclusively or is unaware of the differences between electricity and gas)

Then you need a conversion logic for your `time.Time`s.
This library does the conversion for you.

## How To Use
See [Go Playground](https://go.dev/play/p/rnaj2E2A9xn) for a minimal working example.

Note that this library only modifies timestamps, that are 06:00 German local time (if we're dealing with Gas) or 00:00 German local time (if we're _not_ dealing with Gas).
It won't shift arbitrary timestamps, so in most cases in your application you don't have to manually check if the conversion shall be applied to specific data constellations but only generally think about whether a `time.Time` is interpreted differently by different systems.

## Code Quality / Production Readiness

- The code has [95%](https://github.com/Hochfrequenz/mako_time_converter/blob/main/.github/workflows/coverage.yml#L24) unit test coverage. ✔️
- The package has only one dependency itself (except for testing frameworks) ✔️:
  - [go-playground/validator](https://github.com/go-playground/validator) ️
- No linter warnings in the `golangci-lint` configuration ✔️

## Implicit Requirements

The package requires your relevant timezone data to be present on the system on which you're using it.
It does _not_ include timezone data itself and will panic if the local timezone data is not found.
Please import the [`time/tzdata`](https://pkg.go.dev/time/tzdata) package from the std library, if necessary.

The package does not include any workarounds to actual timezone data (e.g. in the case of Germany calculating the last Sunday in March or October.)
You can do it but you probably shouldn't.
