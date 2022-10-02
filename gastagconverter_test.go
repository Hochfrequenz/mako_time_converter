package mako_time_converter_test

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/hochfrequenz/mako_time_converter"
	"github.com/hochfrequenz/mako_time_converter/enddatetimekind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
}

// SetupSuite sets up the tests
func (s *Suite) SetupSuite() {
}

func (s *Suite) AfterTest(_, _ string) {
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func getBerlinConverter() mako_time_converter.GasTagConverter {
	return mako_time_converter.NewGasTagConverter("Europe/Berlin")
}

func (s *Suite) Test_IsGermanMidnight_true() {
	germanMidnights := []time.Time{
		time.Date(2022, 12, 31, 23, 0, 0, 0, time.UTC),
		time.Date(2022, 06, 15, 22, 0, 0, 0, time.UTC),
		time.Date(2022, 12, 15, 23, 0, 0, 0, time.UTC),
	}
	converter := getBerlinConverter()
	for _, germanMidnight := range germanMidnights {
		then.AssertThat(s.T(), converter.IsGermanMidnight(germanMidnight), is.True())
	}
}

func (s *Suite) Test_IsGermanMidnight_false() {
	notGermanMidnights := []time.Time{
		time.Date(2022, 12, 31, 22, 0, 0, 0, time.UTC),
		time.Date(2022, 06, 15, 23, 0, 0, 0, time.UTC),
		time.Date(2022, 12, 15, 22, 0, 0, 0, time.UTC),
	}
	converter := getBerlinConverter()
	for _, germanMidnight := range notGermanMidnights {
		then.AssertThat(s.T(), converter.IsGermanMidnight(germanMidnight), is.False())
	}
}

func (s *Suite) Test_IsGerman6am_true() {
	german6Ams := []time.Time{
		time.Date(2022, 1, 1, 5, 0, 0, 0, time.UTC),
		time.Date(2022, 06, 15, 4, 0, 0, 0, time.UTC),
		time.Date(2022, 12, 15, 5, 0, 0, 0, time.UTC),
	}
	converter := getBerlinConverter()
	for _, germanMidnight := range german6Ams {
		then.AssertThat(s.T(), converter.IsGerman6Am(germanMidnight), is.True())
	}
}

func (s *Suite) Test_IsGerman6am_false() {
	notGerman6Ams := []time.Time{
		time.Date(2022, 1, 1, 17, 0, 0, 0, time.UTC),
		time.Date(2022, 06, 15, 12, 0, 0, 0, time.UTC),
		time.Date(2022, 12, 15, 23, 0, 0, 0, time.UTC),
	}
	converter := getBerlinConverter()
	for _, germanMidnight := range notGerman6Ams {
		then.AssertThat(s.T(), converter.IsGerman6Am(germanMidnight), is.False())
	}
}

func (s *Suite) Test_German_6Am_To_Midnight_Conversion() {
	pairs := map[time.Time]time.Time{
		time.Date(2022, 12, 31, 5, 0, 0, 0, time.UTC): time.Date(2022, 12, 30, 23, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 1, 5, 0, 0, 0, time.UTC):   time.Date(2022, 12, 31, 23, 0, 0, 0, time.UTC),
		time.Date(2023, 6, 1, 4, 0, 0, 0, time.UTC):   time.Date(2023, 5, 31, 22, 0, 0, 0, time.UTC),
		time.Date(2023, 6, 2, 4, 0, 0, 0, time.UTC):   time.Date(2023, 6, 1, 22, 0, 0, 0, time.UTC),
		time.Date(2023, 03, 26, 4, 0, 0, 0, time.UTC): time.Date(2023, 3, 25, 23, 0, 0, 0, time.UTC),
		time.Date(2023, 10, 29, 5, 0, 0, 0, time.UTC): time.Date(2023, 10, 28, 22, 0, 0, 0, time.UTC),
	}
	converter := getBerlinConverter()
	for german6Am, expectedMidnight := range pairs {
		actualMidnight, err := converter.Convert6AamToMidnight(german6Am)
		then.AssertThat(s.T(), err, is.Nil())
		then.AssertThat(s.T(), actualMidnight, is.EqualTo(expectedMidnight))
		actual6Am, err := converter.ConvertMidnightTo6Am(expectedMidnight)
		then.AssertThat(s.T(), err, is.Nil())
		then.AssertThat(s.T(), actual6Am, is.EqualTo(german6Am))
	}
}

func (s *Suite) Test_German_6Am_To_Midnight_Conversion_Error() {
	converter := getBerlinConverter()
	not6am := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := converter.Convert6AamToMidnight(not6am)
	then.AssertThat(s.T(), err, is.Not(is.Nil()))
}

func (s *Suite) Test_German_Midnight_To_6am_Conversion_Error() {
	converter := getBerlinConverter()
	notGermanMidnight := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := converter.ConvertMidnightTo6Am(notGermanMidnight)
	then.AssertThat(s.T(), err, is.Not(is.Nil()))
}

func (s *Suite) Test_Strip_Time() {
	pairs := map[time.Time]time.Time{
		time.Date(2022, 12, 31, 5, 0, 0, 0, time.UTC): time.Date(2022, 12, 30, 23, 0, 0, 0, time.UTC),
		time.Date(2022, 12, 31, 5, 2, 1, 0, time.UTC): time.Date(2022, 12, 30, 23, 0, 0, 0, time.UTC),
	}
	converter := getBerlinConverter()
	for input, expected := range pairs {
		actual := converter.StripTime(input)
		then.AssertThat(s.T(), actual, is.EqualTo(expected))
	}
}

func pointer[T any](b T) *T {
	return &b
}

func (s *Suite) Test_Gastag_Aware_To_Non_Gastag_Aware() {
	pairs := map[time.Time]time.Time{
		time.Date(2023, 6, 1, 4, 0, 0, 0, time.UTC):  time.Date(2023, 5, 31, 22, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 1, 5, 0, 0, 0, time.UTC): time.Date(2023, 11, 30, 23, 0, 0, 0, time.UTC),
	}
	conversion := mako_time_converter.DateTimeConversionConfiguration{
		Source: mako_time_converter.DateTimeConfiguration{IsGas: true, IsGasTagAware: pointer(true)},
		Target: mako_time_converter.DateTimeConfiguration{IsGas: true, IsGasTagAware: pointer(false)},
	}
	converter := getBerlinConverter()
	for input, expected := range pairs {
		actual, err := converter.Convert(input, conversion)
		then.AssertThat(s.T(), err, is.Nil())
		then.AssertThat(s.T(), actual, is.EqualTo(expected))
	}
	invertedConfig := conversion.Invert()
	for expected, input := range pairs {
		actual, err := converter.Convert(input, invertedConfig)
		then.AssertThat(s.T(), err, is.Nil())
		then.AssertThat(s.T(), actual, is.EqualTo(expected))
	}
}

func (s *Suite) Test_Strom_Inclusive_End_To_Strom_Exclusive_End() {
	pairs := map[time.Time]time.Time{
		time.Date(2023, 05, 30, 22, 0, 0, 0, time.UTC): time.Date(2023, 05, 31, 22, 0, 0, 0, time.UTC),
		time.Date(2023, 05, 31, 22, 0, 0, 0, time.UTC): time.Date(2023, 06, 01, 22, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 31, 23, 0, 0, 0, time.UTC): time.Date(2024, 01, 01, 23, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 01, 23, 0, 0, 0, time.UTC): time.Date(2023, 12, 02, 23, 0, 0, 0, time.UTC),
	}
	conversion := mako_time_converter.DateTimeConversionConfiguration{
		Source: mako_time_converter.DateTimeConfiguration{IsGas: false, IsEndDate: true, EndDateTimeKind: pointer(enddatetimekind.INCLUSIVE)},
		Target: mako_time_converter.DateTimeConfiguration{IsGas: false, IsEndDate: true, EndDateTimeKind: pointer(enddatetimekind.EXCLUSIVE)},
	}
	converter := getBerlinConverter()
	for input, expected := range pairs {
		actual, err := converter.Convert(input, conversion)
		then.AssertThat(s.T(), err, is.Nil())
		then.AssertThat(s.T(), actual, is.EqualTo(expected))
	}
	invertedConfig := conversion.Invert()
	for expected, input := range pairs {
		actual, err := converter.Convert(input, invertedConfig)
		then.AssertThat(s.T(), err, is.Nil())
		then.AssertThat(s.T(), actual, is.EqualTo(expected))
	}
}

func (s *Suite) Test_Gas_Inclusive_End_To_Gas_Exclusive_End() {
	pairs := map[time.Time]time.Time{
		time.Date(2023, 05, 30, 04, 0, 0, 0, time.UTC): time.Date(2023, 05, 31, 04, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 30, 05, 0, 0, 0, time.UTC): time.Date(2023, 12, 31, 05, 0, 0, 0, time.UTC),
		time.Date(2023, 12, 01, 05, 0, 0, 0, time.UTC): time.Date(2023, 12, 02, 05, 0, 0, 0, time.UTC),
	}
	conversion := mako_time_converter.DateTimeConversionConfiguration{
		Source: mako_time_converter.DateTimeConfiguration{IsGas: true, IsGasTagAware: pointer(true), IsEndDate: true, EndDateTimeKind: pointer(enddatetimekind.INCLUSIVE)},
		Target: mako_time_converter.DateTimeConfiguration{IsGas: true, IsGasTagAware: pointer(true), IsEndDate: true, EndDateTimeKind: pointer(enddatetimekind.EXCLUSIVE)},
	}
	converter := getBerlinConverter()
	for input, expected := range pairs {
		actual, err := converter.Convert(input, conversion)
		then.AssertThat(s.T(), err, is.Nil())
		then.AssertThat(s.T(), actual, is.EqualTo(expected))
	}
	invertedConfig := conversion.Invert()
	for expected, input := range pairs {
		actual, err := converter.Convert(input, invertedConfig)
		then.AssertThat(s.T(), err, is.Nil())
		then.AssertThat(s.T(), actual, is.EqualTo(expected))
	}
}

func (s *Suite) Test_Invalid_Configurations_Are_Rejected() {
	invalidConfigs := []mako_time_converter.DateTimeConversionConfiguration{
		{
			Source: mako_time_converter.DateTimeConfiguration{IsGas: false, IsGasTagAware: pointer(true)},
			Target: mako_time_converter.DateTimeConfiguration{IsGas: true, IsGasTagAware: pointer(true)},
		},
		{
			Source: mako_time_converter.DateTimeConfiguration{IsEndDate: true}, // no enddatetime kind given
			Target: mako_time_converter.DateTimeConfiguration{IsEndDate: true, EndDateTimeKind: pointer(enddatetimekind.EXCLUSIVE)},
		},
		{
			Source: mako_time_converter.DateTimeConfiguration{IsGas: true, IsGasTagAware: pointer(true)},
			Target: mako_time_converter.DateTimeConfiguration{IsGas: true}, // no gastag awareness given
		},
	}
	converter := getBerlinConverter()
	for _, invalidConfig := range invalidConfigs {
		_, err := converter.Convert(time.Time{}, invalidConfig)
		then.AssertThat(s.T(), err, is.Not(is.Nil()))
	}
}

func (s *Suite) Test_NewGastagConverter_Panics_for_Unknown_Timezone() {
	assert.Panics(s.T(), func() { mako_time_converter.NewGasTagConverter("OtherContinent/IDontKnow") })
}

func (s *Suite) Test_StripTime_On_Source_Side() {
	pairs := map[time.Time]time.Time{
		time.Date(2023, 05, 30, 22, 1, 2, 3, time.UTC): time.Date(2023, 05, 31, 22, 0, 0, 0, time.UTC)}
	conversion := mako_time_converter.DateTimeConversionConfiguration{
		Source: mako_time_converter.DateTimeConfiguration{IsGas: false, IsEndDate: true, EndDateTimeKind: pointer(enddatetimekind.INCLUSIVE), StripTime: true},
		Target: mako_time_converter.DateTimeConfiguration{IsGas: false, IsEndDate: true, EndDateTimeKind: pointer(enddatetimekind.EXCLUSIVE)},
	}
	converter := getBerlinConverter()
	for input, expected := range pairs {
		actual, err := converter.Convert(input, conversion)
		then.AssertThat(s.T(), err, is.Nil())
		then.AssertThat(s.T(), actual, is.EqualTo(expected))
	}
}

func (s *Suite) Test_StripTime_On_Target_Side() {
	pairs := map[time.Time]time.Time{
		time.Date(2023, 05, 30, 22, 0, 0, 0, time.UTC): time.Date(2023, 05, 30, 22, 0, 0, 0, time.UTC)}
	conversion := mako_time_converter.DateTimeConversionConfiguration{
		Source: mako_time_converter.DateTimeConfiguration{IsGas: true, IsGasTagAware: pointer(false)},
		Target: mako_time_converter.DateTimeConfiguration{IsGas: true, IsGasTagAware: pointer(true), StripTime: true}, // but we'll loose the gastag on target side because of strip time
	}
	converter := getBerlinConverter()
	for input, expected := range pairs {
		actual, err := converter.Convert(input, conversion)
		then.AssertThat(s.T(), err, is.Nil())
		then.AssertThat(s.T(), actual, is.EqualTo(expected))
	}
}

func (s *Suite) Test_Same_Source_And_Target_Leads_To_Utc_Conversion_Only() {
	berlin, _ := time.LoadLocation("Europe/Berlin")
	pairs := map[time.Time]time.Time{
		time.Date(2023, 05, 30, 4, 5, 6, 0, berlin): time.Date(2023, 05, 30, 2, 5, 6, 0, time.UTC)}
	conversion := mako_time_converter.DateTimeConversionConfiguration{
		Source: mako_time_converter.DateTimeConfiguration{},
		Target: mako_time_converter.DateTimeConfiguration{},
	}
	converter := getBerlinConverter()
	for input, expected := range pairs {
		actual, err := converter.Convert(input, conversion)
		then.AssertThat(s.T(), err, is.Nil())
		then.AssertThat(s.T(), actual, is.EqualTo(expected))
	}
}
