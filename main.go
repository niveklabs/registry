package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var version = "dev"

type DiscoveryResponce struct {
	Providers string `json:"providers.v1,omitempty"`
	Modules   string `json:"modules.v1,omitempty"`
	Login     string `json:"login.v1,omitempty"`
}

type ModuleVersions struct {
	Modules []*ModuleProviderVersions `json:"modules"`
}

type ModuleProviderVersions struct {
	Versions []*ModuleVersion `json:"versions"`
}

type ModuleVersion struct {
	Version string `json:"version"`
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Terraform-Get", fmt.Sprintf("https://api.github.com/repos/terraform-aws-modules/terraform-aws-vpc/tarball/v%s?archive=tar.gz", strings.Split(r.RequestURI, "/")[6]))
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w)
}

func versionsHandler(w http.ResponseWriter, r *http.Request) {

	moduleVersions := ModuleVersions{
		Modules: make([]*ModuleProviderVersions, 0),
	}

	moduleVersions.Modules = append(moduleVersions.Modules, &ModuleProviderVersions{Versions: make([]*ModuleVersion, 0)})
	moduleVersions.Modules[0].Versions = append(moduleVersions.Modules[0].Versions, &ModuleVersion{Version: "2.42.0"})

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(moduleVersions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func discoveryHandler(w http.ResponseWriter, r *http.Request) {
	payload := &DiscoveryResponce{
		Modules: "/v1/modules/",
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {

	var wait time.Duration
	var port int
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.IntVar(&port, "port", 8080, "the port the server will listen on")
	flag.Parse()

	r := mux.NewRouter()

	// Add your routes as needed
	r.HandleFunc("/.well-known/terraform.json", discoveryHandler)
	r.HandleFunc("/v1/modules/{namespace}/{name}/{provider}/versions", versionsHandler)
	r.HandleFunc("/v1/modules/{namespace}/{name}/{provider}/{version}/download", downloadHandler)

	r.Use(loggingMiddleware)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Println("Running version", version)
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}

}
