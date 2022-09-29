package conversion

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	mtcconfig "github.com/hochfrequenz/mako_time_converter/configuration"
	"github.com/hochfrequenz/mako_time_converter/configuration/enddatetimekind"
	"log"
	"time"
)

// GasTagConverter is a struct to convert to and from German "Gas-Tag" (which always starts at 6AM German local time)
type GasTagConverter interface {
	// IsGermanMidnight returns true iff the given timestamp is the beginning of a German Stromtag (midnight local time)
	IsGermanMidnight(timestamp time.Time) bool
	// IsGerman6Am returns true if the given timestamp is the beginning of a German Gastag (6AM local time)
	IsGerman6Am(timestamp time.Time) bool
	// Convert6AamToMidnight converts the given local 6Am timestamp to German midnight of the same German day
	Convert6AamToMidnight(timestamp time.Time) (time.Time, error)
	// ConvertMidnightTo6Am converts the given local 6Am timestamp to German midnight of the same German day
	ConvertMidnightTo6Am(timestamp time.Time) (time.Time, error)
	// StripTime removes all hours, minutes, seconds, milliseconds (in german local time) from the given timestamp. This is similar to a "round down" or "floor" in German local time.
	StripTime(timestamp time.Time) time.Time
	// Convert  converts the given timestamp to a DateTimeConversionConfiguration.Target by applying all transformations which are derived from the given configuration time is described by DateTimeConversionConfiguration.Source.
	Convert(timestamp time.Time, configuration mtcconfig.DateTimeConversionConfiguration) (time.Time, error)
}

type locationBasedGasTagConverter struct {
	location *time.Location
}

// NewGasTagConverter returns a GasTagConverter that internally uses the timezone data from the timezone with the given zoneName (e.g. "Europe/Berlin"). It requires the tzdata to be available on the system and will panic if this is not the case.
func NewGasTagConverter(zoneName string) GasTagConverter {
	location, err := time.LoadLocation(zoneName)
	if err != nil {
		errorMsg := fmt.Errorf("The timezone data for '%s' could not be found. Import \"time/tzdata\" anywhere in your project or build with `-tags timetzdata`: https://pkg.go.dev/time/tzdata", zoneName)
		log.Panic(errorMsg)
	}
	return locationBasedGasTagConverter{location: location}
}

// ToLocalTimeConverter contains a method to convert a time into a local time. This will, in most cases, happen on the basis of timezone data, but you are free to write your own conversion, although you're probably missing out on details at one point.
type ToLocalTimeConverter interface {
	// toLocalTime converts a timestamp to a "local" time by adjusting date, time and UTC-offset. The actual point in time in UTC or Unix does _not_ change.
	toLocalTime(timestamp time.Time) time.Time
}

func (l locationBasedGasTagConverter) toLocalTime(timestamp time.Time) time.Time {
	return timestamp.In(l.location)
}

func (l locationBasedGasTagConverter) IsGermanMidnight(timestamp time.Time) bool {
	localTime := l.toLocalTime(timestamp)
	hour, minute, sec := localTime.Clock()
	return hour == 0 && minute == 0 && sec == 0
}

func (l locationBasedGasTagConverter) IsGerman6Am(timestamp time.Time) bool {
	localTime := l.toLocalTime(timestamp)
	hour, minute, sec := localTime.Clock()
	return hour == 6 && minute == 0 && sec == 0
}

func (l locationBasedGasTagConverter) Convert6AamToMidnight(timestamp time.Time) (time.Time, error) {
	if !l.IsGerman6Am(timestamp) {
		return time.Time{}, fmt.Errorf("The given time %v is not German 6am", timestamp)
	}
	return l.StripTime(timestamp), nil
}

func (l locationBasedGasTagConverter) ConvertMidnightTo6Am(timestamp time.Time) (time.Time, error) {
	if !l.IsGermanMidnight(timestamp) {
		return time.Time{}, fmt.Errorf("The given time %v is not German midnight", timestamp)
	}
	localMidnight := l.toLocalTime(timestamp)
	year, month, day := localMidnight.Date()
	local6Am := time.Date(year, month, day, 6, 0, 0, 0, l.location)
	return local6Am.UTC(), nil
}

func (l locationBasedGasTagConverter) StripTime(timestamp time.Time) time.Time {
	localTime := l.toLocalTime(timestamp)
	year, month, day := localTime.Date()
	localMidnight := time.Date(year, month, day, 0, 0, 0, 0, l.location)
	return localMidnight.UTC()
}

func (l locationBasedGasTagConverter) addGermanDay(timestamp time.Time) time.Time {
	localtime := l.toLocalTime(timestamp)
	return localtime.AddDate(0, 0, 1).UTC()
}
func (l locationBasedGasTagConverter) subtractGermanDay(timestamp time.Time) time.Time {
	localtime := l.toLocalTime(timestamp)
	return localtime.AddDate(0, 0, -1).UTC()
}

func (l locationBasedGasTagConverter) Convert(timestamp time.Time, configuration mtcconfig.DateTimeConversionConfiguration) (time.Time, error) {
	validate := validator.New()
	validate.RegisterStructValidation(mtcconfig.DateTimeConversionConfigurationStructLevelValidator, mtcconfig.DateTimeConversionConfiguration{})
	err := validate.Struct(configuration)
	if err != nil {
		return time.Time{}, err
	}
	result := timestamp
	if configuration.Source.StripTime {
		result = l.StripTime(result)
	}
	if configuration.Source == configuration.Target {
		// both are the same, no conversion needed
		return result.UTC(), nil
	}
	if configuration.Source.IsGas { // this implies that the target is also gas, because otherwise the configuration would be invalid
		// handle gas stuff here
		if *configuration.Source.IsGasTagAware && !*configuration.Target.IsGasTagAware {
			// convert from gas-tag to non-gas-tag
			if l.IsGerman6Am(result) {
				result, err = l.Convert6AamToMidnight(result)
				if err != nil {
					return time.Time{}, err
				}
			}
		}
		if !*configuration.Source.IsGasTagAware && *configuration.Target.IsGasTagAware {
			if l.IsGermanMidnight(result) {
				result, err = l.ConvertMidnightTo6Am(result)
				if err != nil {
					return time.Time{}, err
				}
			}
		}
	}
	// else { handle strom-only stuff here }

	if configuration.Source.IsEndDate && configuration.Target.IsEndDate && *configuration.Source.EndDateTimeKind != *configuration.Target.EndDateTimeKind {
		if *configuration.Source.EndDateTimeKind == enddatetimekind.INCLUSIVE { // implicit: target is exclusive
			// convert from inclusive to exclusive
			result = l.addGermanDay(result)
		}
		if *configuration.Source.EndDateTimeKind == enddatetimekind.EXCLUSIVE { // implicit: target is inclusive
			// convert from exclusive to inclusive
			result = l.subtractGermanDay(result)
		}
	}

	if configuration.Target.StripTime {
		result = l.StripTime(result)
	}
	return result.UTC(), nil
}
