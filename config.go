package mako_time_converter

import (
	"github.com/go-playground/validator/v10"
	"github.com/hochfrequenz/mako_time_converter/enddatetimekind"
)

// DateTimeConfiguration describes how a time.Time is meant/interpreted by a system. Two of these configurations allow to convert a time.Time smoothly.
type DateTimeConfiguration struct {
	// IsEndDate is true if the datetime describes an "end date", e.g. a contract end date
	IsEndDate bool `json:"isEndDate"`
	// EndDateTimeKind describes how an end datetime shall be understood (must be set if IsEndDate is true)
	EndDateTimeKind *enddatetimekind.EndDateTimeKind `json:"endDateTimeKind,omitempty" validate:"required_if=IsEndDate true"`
	// IsGas true iff the datetime describes a datetime in Sparte Gas. Please note that this is independent of the information whether the datetime is actually IsGasTagAware. There are systems that discriminate Gas and non-Gas (this is what this flag is for) but are still unaware of the German Gas-Tag.
	IsGas bool `json:"isGas"`
	// IsGasTagAware must be set iff IsGas is true and the date time is aware of the German "Gas-Tag" (meaning that start dates are 6:00 German local time and end dates are 06:00 German local time (if the end date is meant exclusive))
	IsGasTagAware *bool `json:"isGasTagAware,omitempty" validate:"required_if=IsGas true"`
	// Set true to remove all hours, minutes, seconds, milliseconds from the respective time.Time. If set in the DateTimeConversionConfiguration.Source the hours, minutes... will be stripped _before_ the conversion. If set in the DateTimeConversionConfiguration.Target the hours, minutes... will be stripped _after_ the conversion.
	StripTime bool `json:"stripTime"`
}

// A DateTimeConversionConfiguration describes which steps are necessary to convert a datetime from a Source to a Target
type DateTimeConversionConfiguration struct {
	// Source is the configuration of the datetime before the conversion
	Source DateTimeConfiguration `json:"source" validate:"required"`
	// Target is the configuration of the datetime after the conversion
	Target DateTimeConfiguration `json:"target" validate:"required"`
}

// Invert returns an inverted configuration (switched source and target)
func (dtcc *DateTimeConversionConfiguration) Invert() DateTimeConversionConfiguration {
	return DateTimeConversionConfiguration{
		Source: dtcc.Target,
		Target: dtcc.Source,
	}
}

func DateTimeConversionConfigurationStructLevelValidator(sl validator.StructLevel) {
	config := sl.Current().Interface().(DateTimeConversionConfiguration)
	if config.Source.IsGas != config.Target.IsGas {
		sl.ReportError(config.Source, "Source/Target.IsGas", "Target", "Source.IsGas==Target.IsGas", "")
	}
}
