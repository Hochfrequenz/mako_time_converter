# mako_time_converter (go)
![Unittest status badge](https://github.com/hochfrequenz/mako_time_converter/workflows/Unittests/badge.svg)
![Coverage status badge](https://github.com/hochfrequenz/mako_time_converter/workflows/coverage/badge.svg)
![Linter status badge](https://github.com/hochfrequenz/mako_time_converter/workflows/golangci-lint/badge.svg)

Go package to convert between German "Gastag" and "Stromtag", between inclusive and exclusive end dates and combinations of all them.
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
See [Go Playground](https://go.dev/play/p/Xsng2cjWU-Z) for a minimal working example.

## Implicit Requirements

The package requires your relevant timezone data to be present on the system on which you're using it.
It does _not_ include timezone data itself and will panic if the local timezone data is not found.
Please import the [`time/tzdata`](https://pkg.go.dev/time/tzdata) package from the std library, if necessary.

The package does not include any workarounds to actual timezone data (e.g. in the case of Germany calculating the last Sunday in March or October.)
You can do it but you probably shouldn't.
