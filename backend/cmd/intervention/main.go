package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/darylhjd/oams/backend/internal/intervention"
)

const (
	functionsCustomHandlerPort = "FUNCTIONS_CUSTOMHANDLER_PORT"
)

const (
	interventionUrl = "/intervention"
)

func main() {
	port, ok := os.LookupEnv(functionsCustomHandlerPort)
	if ok {
		log.Printf("Custom handler port: %s", port)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(interventionUrl, interventionHandler)

	log.Println("server listening on port ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}

// Response for Azure Function.
type Response struct {
	Outputs     map[string]interface{}
	ReturnValue interface{}
	Logs        []string
}

// interventionHandler handles the invocation of the Early Intervention Service
func interventionHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	service, err := intervention.New(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err = service.Stop(); err != nil {
			log.Fatalf("%s - could not gracefully stop service: %s", intervention.Namespace, err)
		}
	}()

	if err = service.Run(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := &Response{
		Logs: []string{fmt.Sprintf("Early Intervention Service successfully run at %s", now.String())},
	}

	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}
