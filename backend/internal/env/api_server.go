package env

import "os"

const (
	apiServerPort          = "API_SERVER_PORT"
	apiServerAzureTenantId = "API_SERVER_AZURE_TENANT_ID"
	apiServerAzureClientId = "API_SERVER_AZURE_CLIENT_ID"
)

// GetAPIServerPort returns the API_SERVER_PORT environment variable.
func GetAPIServerPort() string {
	return os.Getenv(apiServerPort)
}

// GetAPIServerAzureTenantID returns the API_SERVER_AZURE_TENANT_ID environment variable.
func GetAPIServerAzureTenantID() string {
	return os.Getenv(apiServerAzureTenantId)
}

// GetAPIServerAzureClientID returns the API_SERVER_AZURE_CLIENT_ID environment variable.
func GetAPIServerAzureClientID() string {
	return os.Getenv(apiServerAzureClientId)
}
