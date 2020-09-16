package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var version = "dev"

var registry ModuleVersions
var discovery = &DiscoveryResponce{
	Modules: "/v1/modules/",
}

type DiscoveryResponce struct {
	Providers string `json:"providers.v1,omitempty"`
	Modules   string `json:"modules.v1,omitempty"`
	Login     string `json:"login.v1,omitempty"`
}

type ModuleVersions struct {
	Modules []*ModuleProviderVersions `json:"modules"`
}

type ModuleProviderVersions struct {
	ID       string           `json:"id"`
	Source   string           `json:"source"`
	Versions []*ModuleVersion `json:"versions"`
}

type ModuleVersion struct {
	Download string `json:"download"`
	Version  string `json:"version"`
	Tag      string `json:"tag"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	download := fmt.Sprintf("https://api.github.com/repos/%s/terraform-%s-%s/tarball/v%s?archive=tar.gz", vars["namespace"], vars["provider"], vars["name"], vars["version"])

	w.Header().Set("X-Terraform-Get", download)
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w)
}

func versionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filtered := make([]*ModuleProviderVersions, 0)

	for _, v := range registry.Modules {
		if v.ID == fmt.Sprintf("%s/%s/%s", vars["namespace"], vars["name"], vars["provider"]) {
			filtered = append(filtered, v)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(&ModuleVersions{filtered})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func modulesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(registry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func discoveryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(discovery)
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
	var config string
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.IntVar(&port, "port", 8080, "the port the server will listen on")
	flag.StringVar(&config, "config", "./registry.json", "the json configuration file")
	flag.Parse()

	f, err := os.Open(config)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	json.NewDecoder(f).Decode(&registry)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	// Add your routes as needed
	r.HandleFunc("/", healthHandler)
	r.HandleFunc("/.well-known/terraform.json", discoveryHandler)
	r.HandleFunc("/v1/modules", modulesHandler)
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
