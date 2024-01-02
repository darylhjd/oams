package env

import "os"

const (
	azureEmailEndpoint      = "AZURE_EMAIL_ENDPOINT"
	azureEmailAccessKey     = "AZURE_EMAIL_ACCESS_KEY"
	azureEmailSenderAddress = "AZURE_EMAIL_SENDER_ADDRESS"
)

// GetAzureEmailEndpoint returns the AZURE_EMAIL_ENDPOINT environment variable.
func GetAzureEmailEndpoint() string {
	return os.Getenv(azureEmailEndpoint)
}

// GetAzureEmailAccessKey returns the AZURE_EMAIL_ACCESS_KEY environment variable.
func GetAzureEmailAccessKey() string {
	return os.Getenv(azureEmailAccessKey)
}

// GetAzureEmailSenderAddress returns the AZURE_EMAIL_SENDER_ADDRESS environment variable.
func GetAzureEmailSenderAddress() string {
	return os.Getenv(azureEmailSenderAddress)
}
