package env

const (
	apiServerPort              = "API_SERVER_PORT"
	apiServerAzureTenantId     = "API_SERVER_AZURE_TENANT_ID"
	apiServerAzureClientId     = "API_SERVER_AZURE_CLIENT_ID"
	apiServerAzureClientSecret = "API_SERVER_AZURE_CLIENT_SECRET"
	apiServerAzureLoginScope   = "API_SERVER_AZURE_LOGIN_SCOPE"
)

// GetAPIServerPort returns the API_SERVER_PORT environment variable.
// Note that this variable is required.
func GetAPIServerPort() (string, error) {
	return getRequiredEnv(apiServerPort)
}

// GetAPIServerAzureTenantID returns the API_SERVER_AZURE_TENANT_ID environment variable.
// Note that this variable is required.
func GetAPIServerAzureTenantID() (string, error) {
	return getRequiredEnv(apiServerAzureTenantId)
}

// GetAPIServerAzureClientID returns the API_SERVER_AZURE_CLIENT_ID environment variable.
// Note that this variable is required.
func GetAPIServerAzureClientID() (string, error) {
	return getRequiredEnv(apiServerAzureClientId)
}

// GetAPIServerAzureClientSecret returns the API_SERVER_AZURE_CLIENT_SECRET environment variable.
// Note that this variable is required.
func GetAPIServerAzureClientSecret() (string, error) {
	return getRequiredEnv(apiServerAzureClientSecret)
}

// GetAPIServerAzureLoginScope returns the API_SERVER_AZURE_LOGIN_SCOPE environment variable.
// Note that this variable is required.
func GetAPIServerAzureLoginScope() (string, error) {
	return getRequiredEnv(apiServerAzureLoginScope)
}
