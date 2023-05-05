package apiserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/database"
	"github.com/darylhjd/oams/backend/env"
	"github.com/darylhjd/oams/backend/logger"
	"github.com/darylhjd/oams/backend/servers/apiserver/v1"
)

const (
	Namespace = "apiserver"

	microsoftAuthority = "https://login.microsoftonline.com/%s/"

	v1Url = "/api/v1/"
)

// APIServer defines the server structure for the OAMS API service.
type APIServer struct {
	l  *zap.Logger
	db *database.DB

	v1 *v1.APIServerV1
}

// Start the APIServer.
func (s *APIServer) Start() error {
	s.l.Info(fmt.Sprintf("%s - starting service...", Namespace))

	port := env.GetAPIServerPort()
	s.l.Info(fmt.Sprintf("%s - service started", Namespace), zap.String("port", port))
	return http.ListenAndServe(fmt.Sprintf(":%s", port), s)
}

// Stop closes any external connections (e.g. database) and stops the server gracefully.
func (s *APIServer) Stop() error {
	return s.db.Close()
}

func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	// To add more versions for this URL, simply add another handler for it.
	mux.Handle(v1Url, http.StripPrefix(strings.TrimSuffix(v1Url, "/"), s.v1))

	mux.ServeHTTP(w, r)
}

func (s *APIServer) GetLogger() *zap.Logger {
	return s.l
}

// New creates a new APIServer. Use Start() to start the server.
func New() (*APIServer, error) {
	l, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("%s - failed to initialise: %w", Namespace, err)
	}

	db, err := database.Connect()
	if err != nil {
		return nil, fmt.Errorf("%s - could not connect to database: %w", Namespace, err)
	}

	cred, err := confidential.NewCredFromSecret(env.GetAPIServerAzureClientSecret())
	if err != nil {
		return nil, fmt.Errorf("%s - could not create credential from client secret: %w", Namespace, err)
	}

	azureClient, err := confidential.New(
		fmt.Sprintf(microsoftAuthority, env.GetAPIServerAzureTenantID()),
		env.GetAPIServerAzureClientID(),
		cred)
	if err != nil {
		return nil, fmt.Errorf("%s - could not create azure client: %w", Namespace, err)
	}

	return &APIServer{l, db, v1.NewAPIServerV1(l, db, &azureClient)}, nil
}
