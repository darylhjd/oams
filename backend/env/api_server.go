package env

import "os"

const (
	apiServerPort                  = "API_SERVER_PORT"
	apiServerAzureTenantId         = "API_SERVER_AZURE_TENANT_ID"
	apiServerAzureClientId         = "API_SERVER_AZURE_CLIENT_ID"
	apiServerAzureClientSecret     = "API_SERVER_AZURE_CLIENT_SECRET"
	apiServerAzureLoginScope       = "API_SERVER_AZURE_LOGIN_SCOPE"
	apiServerAzureLoginCallbackUrl = "API_SERVER_AZURE_LOGIN_CALLBACK_URL"
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

// GetAPIServerAzureClientSecret returns the API_SERVER_AZURE_CLIENT_SECRET environment variable.
func GetAPIServerAzureClientSecret() string {
	return os.Getenv(apiServerAzureClientSecret)
}

// GetAPIServerAzureLoginScope returns the API_SERVER_AZURE_LOGIN_SCOPE environment variable.
func GetAPIServerAzureLoginScope() string {
	return os.Getenv(apiServerAzureLoginScope)
}

// GetAPIServerAzureLoginCallbackURL returns the API_SERVER_AZURE_LOGIN_CALLBACK_URL environment variable.
func GetAPIServerAzureLoginCallbackURL() string {
	return os.Getenv(apiServerAzureLoginCallbackUrl)
}
