package microservicetransport

import "net/http"

// ServiceTransport - Interface responsible for communication.
type ServiceTransport interface {
	// Call - Do the current service request.
	//
	// Return:
	//     *http.Response - The response of the request.
	//     error - An error if it occurred.
	Call() (*http.Response, error)

	// Dial - Create a request to a service resource.
	//
	// Params:
	//     request *Request - The request to perform against the service.
	//
	// Return:
	//     error - An error if it occurred.
	Dial(request *Request) error
}
