package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	//"github.com/darylhjd/oams/backend/internal/intervention"
)

// InvokeResponse ...
type InvokeResponse struct {
	Outputs     map[string]interface{}
	ReturnValue interface{}
	Logs        []string
}

func interventionHandler(w http.ResponseWriter, r *http.Request) {
	// service, err := intervention.New(r.Context())
	// if err != nil {
	// 	log.Fatalf("%s - cannot start service: %s", intervention.Namespace, err)
	// }
	// defer func() {
	// 	if err = service.Stop(); err != nil {
	// 		log.Printf("%s - could not gracefully stop service: %s", intervention.Namespace, err)
	// 	}
	// }()

	// if err = service.Run(r.Context()); err != nil {
	// 	log.Printf("%s - error running service: %s", intervention.Namespace, err)
	// }
	defer func() { _ = r.Body.Close() }()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	response := &InvokeResponse{
		ReturnValue: fmt.Sprintf("%s", data),
		Outputs: map[string]interface{}{
			"result": true,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func main() {
	port, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !ok {
		port = "3000"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/intervention", interventionHandler)

	log.Println("function server listening on port ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
