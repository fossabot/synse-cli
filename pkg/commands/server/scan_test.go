package server

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gotestyourself/gotestyourself/assert"
	"github.com/gotestyourself/gotestyourself/golden"

	"github.com/vapor-ware/synse-cli/internal/test"
	"github.com/vapor-ware/synse-cli/pkg/config"
)

const (
	// the mocked 200 OK JSON response for the Synse Server 'scan' route
	scanRespOK = `
{
  "racks":[
    {
      "id":"rack-1",
      "boards":[
        {
          "id":"board-1",
          "devices":[
            {
              "id":"device-1",
              "info":"Synse Temperature Sensor",
              "type":"temperature"
            },
            {
              "id":"device-2",
              "info":"Synse Fan",
              "type":"fan"
            },
            {
              "id":"device-3",
              "info":"Synse LED",
              "type":"led"
            },
            {
              "id":"device-4",
              "info":"Synse LED",
              "type":"led"
            },
            {
              "id":"device-5",
              "info":"Synse Temperature Sensor",
              "type":"temperature"
            },
            {
              "id":"device-6",
              "info":"Synse Temperature Sensor",
              "type":"temperature"
            },
            {
              "id":"device-7",
              "info":"Synse Temperature Sensor",
              "type":"temperature"
            },
            {
              "id":"device-8",
              "info":"Synse Temperature Sensor",
              "type":"temperature"
            }
          ]
        }
      ]
    }
  ]
}`

	// the mocked 500 error JSON response for the Synse Server 'scan' route
	scanRespErr = `
{
  "http_code":500,
  "error_id":0,
  "description":"unknown",
  "timestamp":"2018-03-14 15:34:42.243715",
  "context":"test error."
}`
)

// TestScanCommandError tests the 'scan' command when it is unable to
// connect to the Synse Server instance because the active host is nil.
func TestScanCommandError(t *testing.T) {
	test.Setup()

	app := test.NewFakeApp()
	app.Commands = append(app.Commands, ServerCommand)

	err := app.Run([]string{
		app.Name,
		ServerCommand.Name,
		scanCommand.Name,
	})

	assert.Assert(t, golden.String(app.ErrBuffer.String(), "error.nil.golden"))
	test.ExpectExitCoderError(t, err)
}

// TestScanCommandError2 tests the 'scan' command when it is unable to
// connect to the Synse Server instance because the active host is not a
// Synse Server instance.
func TestScanCommandError2(t *testing.T) {
	test.Setup()
	config.Config.ActiveHost = &config.HostConfig{
		Name:    "test-host",
		Address: "localhost:5151",
	}

	app := test.NewFakeApp()
	app.Commands = append(app.Commands, ServerCommand)

	err := app.Run([]string{
		app.Name,
		ServerCommand.Name,
		scanCommand.Name,
	})

	assert.Assert(t, golden.String(app.ErrBuffer.String(), "error.bad_host.golden"))
	test.ExpectExitCoderError(t, err)
}

// TestScanCommandRequestError tests the 'scan' command when it gets a
// 500 response from Synse Server.
func TestScanCommandRequestError(t *testing.T) {
	test.Setup()

	mux, server := test.Server()
	defer server.Close()
	mux.HandleFunc(
		"/synse/2.0/scan",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			fmt.Fprint(w, scanRespErr)
		},
	)

	test.AddServerHost(server)
	app := test.NewFakeApp()
	app.Commands = append(app.Commands, ServerCommand)

	err := app.Run([]string{
		app.Name,
		ServerCommand.Name,
		scanCommand.Name,
	})

	assert.Assert(t, golden.String(app.ErrBuffer.String(), "error.500.golden"))
	test.ExpectExitCoderError(t, err)
}

// TestScanCommandRequestErrorYaml tests the 'scan' command when it gets
// a 200 response from Synse Server, with YAML output.
func TestScanCommandRequestErrorYaml(t *testing.T) {
	test.Setup()

	mux, server := test.Server()
	defer server.Close()
	mux.HandleFunc(
		"/synse/2.0/scan",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, scanRespOK)
		},
	)

	test.AddServerHost(server)
	app := test.NewFakeApp()
	app.Commands = append(app.Commands, ServerCommand)

	err := app.Run([]string{
		app.Name,
		"--format", "yaml",
		ServerCommand.Name,
		scanCommand.Name,
	})

	assert.Assert(t, golden.String(app.ErrBuffer.String(), "scan.error.yaml.golden"))
	test.ExpectExitCoderError(t, err)
}

// TestScanCommandRequestErrorJson tests the 'scan' command when it gets
// a 200 response from Synse Server, with JSON output.
func TestScanCommandRequestErrorJson(t *testing.T) {
	test.Setup()

	mux, server := test.Server()
	defer server.Close()
	mux.HandleFunc(
		"/synse/2.0/scan",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, scanRespOK)
		},
	)

	test.AddServerHost(server)
	app := test.NewFakeApp()
	app.Commands = append(app.Commands, ServerCommand)

	err := app.Run([]string{
		app.Name,
		"--format", "json",
		ServerCommand.Name,
		scanCommand.Name,
	})

	assert.Assert(t, golden.String(app.ErrBuffer.String(), "scan.error.json.golden"))
	test.ExpectExitCoderError(t, err)
}

// TestScanCommandRequestSuccessPretty tests the 'scan' command when it gets
// a 200 response from Synse Server, with pretty output.
func TestScanCommandRequestSuccessPretty(t *testing.T) {
	test.Setup()

	mux, server := test.Server()
	defer server.Close()
	mux.HandleFunc(
		"/synse/2.0/scan",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, scanRespOK)
		},
	)

	test.AddServerHost(server)
	app := test.NewFakeApp()
	app.Commands = append(app.Commands, ServerCommand)

	err := app.Run([]string{
		app.Name,
		"--format", "pretty",
		ServerCommand.Name,
		scanCommand.Name,
	})

	assert.Assert(t, golden.String(app.OutBuffer.String(), "scan.success.pretty.golden"))
	test.ExpectNoError(t, err)
}
