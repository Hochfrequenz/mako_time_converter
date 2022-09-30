package mako_time_converter_test

import (
	"encoding/json"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/hochfrequenz/mako_time_converter"
	"github.com/hochfrequenz/mako_time_converter/enddatetimekind"
	"strings"
)

func (s *Suite) Test_Configuration_Serialization() {
	configs := []mako_time_converter.DateTimeConversionConfiguration{
		{Source: mako_time_converter.DateTimeConfiguration{IsGas: true, IsGasTagAware: pointer(true), IsEndDate: true, EndDateTimeKind: pointer(enddatetimekind.EXCLUSIVE)},
			Target: mako_time_converter.DateTimeConfiguration{IsGas: true, IsGasTagAware: pointer(false), IsEndDate: true, EndDateTimeKind: pointer(enddatetimekind.INCLUSIVE)}},
	}
	for _, config := range configs {
		jsonBytes, err := json.Marshal(config)
		then.AssertThat(s.T(), err, is.Nil())
		then.AssertThat(s.T(), strings.Contains(string(jsonBytes), "EXCLUSIVE"), is.True()) // json string enum serialization
		var deserializedConfig mako_time_converter.DateTimeConversionConfiguration
		err = json.Unmarshal(jsonBytes, &deserializedConfig)
		then.AssertThat(s.T(), err, is.Nil())
		then.AssertThat(s.T(), config, is.EqualTo(deserializedConfig))
	}
}
