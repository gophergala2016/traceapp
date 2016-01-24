package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/common/log"
)

func main() {
	r := httprouter.New()

	// r.GET("/overview", Overview)
	r.GET("/graph/*importpath", Grapher)

	r.ServeFiles("/static/*filepath", http.Dir("static/"))

	log.Infoln("launching tracegraph on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
