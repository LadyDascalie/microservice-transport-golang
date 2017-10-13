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
	Protocol       string
}

// getProtocol - Get the transfer protocol to use for the service
func (s *Service) getProtocol() string {
	switch s.Protocol {
	case config.ProtocolHTTP, config.ProtocolHTTPS:
		return s.Protocol
	default:
		return config.ProtocolHTTPS
	}
}

// Call - Do the current service request.
func (s *Service) Call() (*http.Response, error) {
	return http.DefaultClient.Do(s.CurrentRequest)
}

// Dial - Create a request to a service resource.
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
	resourceUrl := fmt.Sprintf("%s://%s/%s", config.ProtocolHTTP, dnsName, request.Resource)

	// Append the query string if we have any.
	if len(request.Query) > 0 {
		resourceUrl = fmt.Sprintf("%s?%s", resourceUrl, request.Query.Encode())
	}

	// Create the request.
	s.CurrentRequest, err = http.NewRequest(request.Method, resourceUrl, request.Body)
	if err != nil {
		return err
	}

	return nil
}
