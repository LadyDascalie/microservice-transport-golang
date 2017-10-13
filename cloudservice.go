package microservicetransport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"strconv"

	"github.com/LUSHDigital/microservice-core-golang/response"
	"github.com/LUSHDigital/microservice-transport-golang/config"
	"github.com/LUSHDigital/microservice-transport-golang/domain"
	"github.com/LUSHDigital/microservice-transport-golang/models"
	"github.com/pkg/errors"
)

// AuthCredentials - Credentials needed to authenticate for a cloud service.
type AuthCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CloudService - Responsible for communication with a cloud service.
type CloudService struct {
	Service
	Credentials *AuthCredentials
}

// authenticate - Authenticate against the API gateway and return an auth token.
func (c *CloudService) authenticate() (*models.Token, error) {
	loginBody := new(bytes.Buffer)
	if err := json.NewEncoder(loginBody).Encode(c.Credentials); err != nil {
		return nil, fmt.Errorf("cannot encode json: %s", err)
	}

	loginReq, err := http.NewRequest(http.MethodPost, c.GetApiGatewayUrl(), loginBody)
	if err != nil {
		return nil, fmt.Errorf("cannot build login request: %s", err)
	}

	loginResp, err := http.DefaultClient.Do(loginReq)
	if err != nil {
		return nil, fmt.Errorf("cannot perform login request: %s", err)
	}

	// Decode response.
	serviceResponse := response.Response{}
	jsonErr := json.NewDecoder(loginResp.Body).Decode(&serviceResponse)
	if jsonErr != nil {
		return nil, fmt.Errorf("cannot decode login response: %s", err)
	}

	// Handle any error codes.
	switch loginResp.StatusCode {
	// Custom error for login failed.
	case http.StatusUnauthorized, http.StatusNotFound:
		return nil, LoginUnauthorisedError{}

	// 200 and 304 are all good.
	case http.StatusOK, http.StatusNotModified:
		break

	// Something somewhere broken!
	default:
		return nil, fmt.Errorf("api gateway login failed: %s", serviceResponse.Message)
	}

	// Extract the consumer from the response.
	var consumer *models.Consumer
	consumerErr := serviceResponse.ExtractData("consumer", &consumer)
	if consumerErr != nil {
		return nil, fmt.Errorf("could not extract consumer data: %s", consumerErr)
	}

	if len(consumer.Tokens) == 0 {
		return nil, ConsumerHasNoTokensError{}
	}

	return consumer.Tokens[0], nil
}

// GetApiGatewayUrl - Get the url of the API gateway.
func (c *CloudService) GetApiGatewayUrl() string {
	if c.Environment == "staging" {
		return fmt.Sprintf("%s-%s.%s", config.GetGatewayUri(), c.Environment, config.GetServiceDomain())
	}

	return fmt.Sprintf("%s.%s", config.GetGatewayUri(), config.GetServiceDomain())
}

// Call - Do the current service request.
func (c *CloudService) Call() (*http.Response, error) {
	return http.DefaultClient.Do(c.CurrentRequest)
}

// Dial - Create a request to a service resource.
func (c *CloudService) Dial(request *Request) error {
	if c.Credentials.Email == "" || c.Credentials.Password == "" {
		return errors.New("cannot authenticate for cloud service: missing credentials")
	}

	token, err := c.authenticate()
	if err != nil {
		return fmt.Errorf("cannot authenticate for cloud service: %s", err)
	}

	// Make any alterations based upon the namespace.
	switch c.Namespace {
	case "aggregators":
		c.Name = strings.Join([]string{config.AggregatorDomainPrefix, c.Name}, "-")
	}

	// Determine the service namespace to use based on the service version.
	serviceNamespace := c.Name
	if c.Version != 0 {
		serviceNamespace = fmt.Sprintf("%s-%d", serviceNamespace, c.Version)
	}

	cloudServiceUrl := domain.BuildCloudServiceUrl(c.getProtocol(), c.GetApiGatewayUrl(), serviceNamespace, c.Name)

	// Build the resource URL.
	resourceUrl := fmt.Sprintf("%s/%s", cloudServiceUrl, request.Resource)

	// Append the query string if we have any.
	if len(request.Query) > 0 {
		resourceUrl = fmt.Sprintf("%s?%s", resourceUrl, request.Query.Encode())
	}

	// Create the request.
	var reqErr error
	c.CurrentRequest, reqErr = http.NewRequest(request.Method, resourceUrl, request.Body)
	if reqErr != nil {
		return reqErr
	}

	// Set the auth token header.
	c.CurrentRequest.Header.Set(config.AuthHeader, token.PrepareForHttp())

	// Add the version header to the request if applicable.
	if c.Version != 0 {
		c.CurrentRequest.Header.Set(config.ServiceVersionHeader, strconv.Itoa(c.Version))
	}

	return nil
}
