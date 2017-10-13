package domain

import "fmt"

// BuildServiceDNSName - Build the full DNS name for a service.
//
// Params:
//     service string - The machine name of the service.
//     branch string - The VCS branch to use for the service.
//     environment string - The CI environment to use for the service.
//     serviceNamespace string - The k8s namespace of the service.
//
// Return:
//     string - The fully qualified DNS name of the service.
//     error - An error if it occurred.
func BuildServiceDNSName(service, branch, environment, serviceNamespace string) string {
	return fmt.Sprintf("%s-%s-%s.%s", service, branch, environment, serviceNamespace)
}
