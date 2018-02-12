package client

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/urfave/cli"
	"github.com/vapor-ware/synse-cli/internal/test"
)

func TestCheck1(t *testing.T) {
	var err error
	resp := http.Response{
		StatusCode: 200,
	}

	e := check(&resp, err)
	test.ExpectNoError(t, e)
}

func TestCheck2(t *testing.T) {
	err := fmt.Errorf("test error")
	resp := http.Response{
		StatusCode: 200,
	}

	e := check(&resp, err)
	test.ExpectExitCoderError(t, e)
}

func TestCheck3(t *testing.T) {
	var err error
	resp := http.Response{
		StatusCode: 500,
	}

	e := check(&resp, err)
	test.ExpectExitCoderError(t, e)
}

func TestMakeURI(t *testing.T) {
	var testTable = []struct {
		in  []string
		out string
	}{
		{
			in:  []string{""},
			out: "",
		},
		{
			in:  []string{"foo"},
			out: "foo",
		},
		{
			in:  []string{"foo", "bar"},
			out: "foo/bar",
		},
		{
			in:  []string{"foo", "bar", "baz/"},
			out: "foo/bar/baz/",
		},
		{
			in:  []string{"foo/", "/bar"},
			out: "foo///bar",
		},
	}

	cli.OsExiter = func(int) {}

	for _, testCase := range testTable {
		r := MakeURI(testCase.in...)
		if r != testCase.out {
			t.Errorf("MakeURI(%v) => %s, want %s", testCase.in, r, testCase.out)
		}
	}
}

func TestDoGet(t *testing.T) {
	test.Setup()
	cli.OsExiter = func(int) {}

	type testData struct {
		Status string
	}

	mux, server := test.Server()
	defer server.Close()
	mux.HandleFunc("/synse/2.0/foobar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": "ok"}`)
	})

	test.AddServerHost(server)

	out := &testData{}
	err := DoGet("foobar", out)

	test.ExpectNoError(t, err)

	if out.Status != "ok" {
		t.Errorf("expected status 'ok', but was %s", out.Status)
	}
}

func TestDoGetUnversioned(t *testing.T) {
	test.Setup()
	cli.OsExiter = func(int) {}

	type testData struct {
		Status string
	}

	mux, server := test.UnversionedServer()
	defer server.Close()
	mux.HandleFunc("/synse/foobar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": "ok"}`)
	})

	test.AddServerHost(server)

	out := &testData{}
	err := DoGet("foobar", out)

	test.ExpectNoError(t, err)

	if out.Status != "ok" {
		t.Errorf("expected status 'ok', but was %s", out.Status)
	}
}

func TestDoPost(t *testing.T) {
	test.Setup()
	cli.OsExiter = func(int) {}

	type testData struct {
		Status string
	}

	mux, server := test.Server()
	defer server.Close()
	mux.HandleFunc("/synse/2.0/foobar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected request to be POST, but was %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": "ok"}`)
	})

	test.AddServerHost(server)

	out := &testData{}
	err := DoPost("foobar", testData{}, out)

	test.ExpectNoError(t, err)

	if out.Status != "ok" {
		t.Errorf("expected status 'ok', but was %s", out.Status)
	}
}