package microservicetransport

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/LUSHDigital/microservice-transport-golang/config"
	"github.com/LUSHDigital/microservice-transport-golang/domain"
)

// Service - Responsible for communication with a service.
type Service struct {
	Branch         string
	CurrentRequest *http.Request
	Environment    string
	Namespace      string
	Name           string
	Version        int
}

// Call - Do the current service request.
//
// Return:
//     *http.Response - The response of the request.
//     error - An error if it occurred.
func (s *Service) Call() (*http.Response, error) {
	return http.DefaultClient.Do(s.CurrentRequest)
}

// Dial - Create a request to a service resource.
//
// Params:
//     method string - The HTTP method to use for the request.
//     resource string - The service resource we want to request.
//     body io.Reader - The body to pass to the request. Can be nil.
//
// Return:
//     error - An error if it occurred.
func (s *Service) Dial(request *Request) error {
	var err error

	// Make any alterations based upon the namespace.
	switch s.Namespace {
	case "aggregators":
		s.Name = strings.Join([]string{config.AggregatorDomainPrefix, s.Name}, "-")
	}

	// Determine the service namespace to use based on the service version.
	serviceNamespace := s.Name
	if s.Version != 0 {
		serviceNamespace = fmt.Sprintf("%s-%d", serviceNamespace, s.Version)
	}

	// Get the name of the service.
	dnsName := domain.BuildServiceDNSName(s.Name, s.Branch, s.Environment, serviceNamespace)

	// Build the resource URL.
	resourceURL := fmt.Sprintf("%s://%s/%s", config.ProtocolHTTP, dnsName, request.Resource)

	// Append the query string if we have any.
	if len(request.Query) > 0 {
		resourceURL = fmt.Sprintf("%s?%s", resourceURL, request.Query.Encode())
	}

	// Create the request.
	s.CurrentRequest, err = http.NewRequest(request.Method, resourceURL, request.Body)
	if err != nil {
		return err
	}

	return nil
}
