package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
)

func main() {
	// viper stuff
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.AutomaticEnv()

	r := httprouter.New()

	// r.NotFound = http.RedirectHandler("/static/404.html", http.StatusMovedPermanently)

	// s := sessions.NewCookieStore([]byte(viper.GetString("SESSION_SECRET")))
	// // add Oauth2 callback/auth URLs
	// githubLogin(r, s)
	//
	// r.GET("/", func(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// 	t, _ := template.ParseFiles("./template/layout.html.tmpl")
	// 	t.Execute(res, nil)
	// })

	// r.GET("/overview", Overview)
	r.GET("/graph/*importpath", Grapher)

	r.ServeFiles("/static/*filepath", http.Dir("static/"))

	log.Infoln("launching traceapp on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
