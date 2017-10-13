package domain

import "fmt"

// BuildServiceDNSName - Build the full DNS name for a service.
func BuildServiceDNSName(service, branch, environment, serviceNamespace string) string {
	return fmt.Sprintf("%s-%s-%s.%s", service, branch, environment, serviceNamespace)
}

// BuildCloudServiceUrl - Build the full URL for a cloud service.
func BuildCloudServiceUrl(apiGatewayUrl, serviceNamespace, serviceName string) string {
	return fmt.Sprintf("%s/%s/%s", apiGatewayUrl, serviceNamespace, serviceName)
}
