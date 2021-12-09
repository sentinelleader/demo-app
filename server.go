package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/integrations/nrgorilla"
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
)

func helloWorldHandler() http.Handler {
	return reqdHandler("Hellow World")
}

func statusHandler() http.Handler {
	return reqdHandler("OK")
}

func failureHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
	})
}

func reqdHandler(msg string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Handling request for %s\n", r.URL.Path)
		w.Write([]byte(msg))
	})
}

func main() {

	licenseKeyFile := flag.String("licenseKeyFile", "", "Path to the licensekey file")
	crashMode := flag.Bool("crashMode", false, "whether app should fail requests")
	flag.Parse()

	licenseKey, err := ioutil.ReadFile(*licenseKeyFile)
	if err != nil {
		panic(err)
	}

	if len(licenseKey) == 0 {
		panic("Empty license key file")
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("Test Gorilla App"),
		newrelic.ConfigLicense(string(licenseKey)),
		newrelic.ConfigInfoLogger(os.Stdout),
	)
	if nil != err {
		panic(err)
	}

	r := mux.NewRouter()
	r.Use(nrgorilla.Middleware(app))

	if *crashMode {
		// fail purposefully
		r.Handle("/", failureHandler())
	} else {
		r.Handle("/", helloWorldHandler())
	}
	r.Handle("/status", statusHandler())

	_, r.NotFoundHandler = newrelic.WrapHandle(app, "NotFoundHandler", reqdHandler("not found"))
	_, r.MethodNotAllowedHandler = newrelic.WrapHandle(app, "MethodNotAllowedHandler", reqdHandler("method not allowed"))

	fmt.Println("Starting webserver ...")
	http.ListenAndServe(":8000", r)
}
