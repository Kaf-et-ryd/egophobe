package server

import (
	"errors"
	"io/ioutil"
	"log"
	"meguca/config"
	"meguca/db"
	"meguca/lang"
	"meguca/templates"
	"meguca/util"
	"net/http"
	"os"
	"testing"

	"github.com/dimfeld/httptreemux"
)

// Global router used for tests
var router http.Handler

func init() {
	isTest = true
	router = createRouter()
	webRoot = "testdata"
	imageWebRoot = "testdata"
	db.ConnArgs = db.TestConnArgs
	db.IsTest = true

	if err := db.LoadDB(); err != nil {
		panic(err)
	}
	config.Set(config.Configs{
		Public: config.Public{
			DefaultLang: "fr_FR",
		},
	})
	if err := util.Waterfall(lang.Load, templates.Compile); err != nil {
		panic(err)
	}
}

func TestAllBoardRedirect(t *testing.T) {
	t.Parallel()

	rec, req := newPair("/")
	router.ServeHTTP(rec, req)
	assertCode(t, rec, 301)

	loc := rec.Header().Get("Location")
	if loc != "/all/" {
		t.Fatalf("unexpected redirect result: %s", loc)
	}
}

func TestPanicHandler(t *testing.T) {
	r := httptreemux.NewContextMux()
	h := func(_ http.ResponseWriter, _ *http.Request) {
		panic(errors.New("foo"))
	}
	r.GET("/panic", h)
	r.PanicHandler = text500
	rec, req := newPair("/panic")

	// Prevent printing stack trace to terminal
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	r.ServeHTTP(rec, req)
	assertCode(t, rec, 500)
	assertBody(t, rec, "500 foo\n")
}

func TestGzip(t *testing.T) {
	enableGzip = true
	defer func() {
		enableGzip = false
	}()

	r := createRouter()
	rec, req := newPair("/json/config")
	req.Header.Set("Accept-Encoding", "gzip")

	r.ServeHTTP(rec, req)

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatal("response not gzipped")
	}
}
