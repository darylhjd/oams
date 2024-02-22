package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/darylhjd/oams/backend/internal/env"
	"github.com/darylhjd/oams/backend/pkg/azmail"
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
func interventionHandler(w http.ResponseWriter, _ *http.Request) {
	now := time.Now()

	//service, err := intervention.New(r.Context())
	//if err != nil {
	//	log.Fatalf("%s - cannot start service: %s", intervention.Namespace, err)
	//}
	//defer func() {
	//	if err = service.Stop(); err != nil {
	//		log.Fatalf("%s - could not gracefully stop service: %s", intervention.Namespace, err)
	//	}
	//}()
	//
	//if err = service.Run(r.Context()); err != nil {
	//	log.Fatalf("%s - error running service: %s", intervention.Namespace, err)
	//}

	mailClient, err := azmail.NewClient(
		env.GetAzureEmailEndpoint(), env.GetAzureEmailAccessKey(), env.GetAzureEmailSenderAddress(),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mail := azmail.NewMail()
	mail.Recipients = azmail.MailRecipients{
		To: []azmail.MailAddress{{
			Address:     "harj0002@e.ntu.edu.sg",
			DisplayName: "Daryl",
		}},
	}
	mail.Content = azmail.MailContent{
		Subject:   "This is a test email from Azure Functions",
		PlainText: "Dummy content.",
		Html:      "Dummy content.",
	}

	if err = mailClient.SendMails(mail); err != nil {
		log.Fatalf("could not send mails: %s", err)
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
