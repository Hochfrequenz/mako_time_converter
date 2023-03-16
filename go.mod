module github.com/hochfrequenz/mako_time_converter

go 1.19

require (
	github.com/corbym/gocrest v1.0.6
	github.com/go-playground/validator/v10 v10.11.2
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/hochfrequenz/mako_time_converter/configuration => ./configuration
	github.com/hochfrequenz/mako_time_converter/conversion => ./conversion
)
