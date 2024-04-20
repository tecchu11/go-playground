package nrmux_test

import (
	"testing"

	"github.com/newrelic/go-agent/v3/newrelic"
)

var app *newrelic.Application

func TestMain(m *testing.M) {
	var err error
	app, err = newrelic.NewApplication(
		newrelic.ConfigLicense("0000000000000000000000000000000000000000"),
		newrelic.ConfigAppName("test-local"),
	)
	if err != nil {
		panic(err)
	}
	m.Run()
}
