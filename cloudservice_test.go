package microservicetransport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"net/http/httptest"

	"os"

	"fmt"

	"github.com/LUSHDigital/microservice-core-golang/format"
	"github.com/LUSHDigital/microservice-core-golang/response"
	"github.com/LUSHDigital/microservice-transport-golang/config"
	"github.com/LUSHDigital/microservice-transport-golang/models"
)

func TestCloudService_Dial(t *testing.T) {
	// Start a HTTP server to act as a fake API gateway.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := response.New(http.StatusOK, response.StatusOk, "", &response.Data{
			Type: "consumer",
			Content: models.Consumer{
				Tokens: []*models.Token{
					{
						Type:  "JWT",
						Value: "xxxx.xxxx.xxxx",
					},
				},
			},
		})
		format.JSONResponseFormatter(w, resp)
	}))
	defer ts.Close()

	// Set the URL of the fake gateway as an environment variable.
	os.Setenv("SOA_GATEWAY_URL", ts.URL)

	tt := []struct {
		name           string
		service        CloudService
		request        *Request
		postData       map[string]string
		expectedMethod string
		expectedUri    string
		expectedBody   string
	}{
		{
			name: "GET HTTP",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
				},
			},
			request: &Request{
				Method:   http.MethodGet,
				Resource: "things",
				Protocol: config.ProtocolHTTP,
			},
			expectedUri: "things",
		},
		{
			name: "GET HTTPS",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
				},
			},
			request: &Request{
				Method:   http.MethodGet,
				Resource: "things",
				Protocol: config.ProtocolHTTPS,
			},
			expectedUri: "things",
		},
		{
			name: "GET with query HTTP",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
				},
			},
			request: &Request{
				Method:   http.MethodGet,
				Resource: "things",
				Protocol: config.ProtocolHTTP,
				Query: url.Values{
					"baz": []string{"qux"},
					"foo": []string{"bar"},
				},
			},
			expectedUri: "things?baz=qux&foo=bar",
		},
		{
			name: "GET with query HTTPS",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
				},
			},
			request: &Request{
				Method:   http.MethodGet,
				Resource: "things",
				Protocol: config.ProtocolHTTPS,
				Query: url.Values{
					"foo": []string{"bar"},
					"baz": []string{"qux"},
				},
			},
			expectedUri: "things?baz=qux&foo=bar",
		},
		{
			name: "POST HTTP",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
				},
			},
			postData: map[string]string{
				"foo": "bar",
				"baz": "qux",
			},
			request: &Request{
				Method:   http.MethodPost,
				Resource: "things",
				Protocol: config.ProtocolHTTP,
			},
			expectedUri: "things",
		},
		{
			name: "POST HTTPS",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
				},
			},
			postData: map[string]string{
				"foo": "bar",
				"baz": "qux",
			},
			request: &Request{
				Method:   http.MethodPost,
				Resource: "things",
				Protocol: config.ProtocolHTTPS,
			},
			expectedUri: "things",
		},
		{
			name: "POST with query HTTP",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
				},
			},
			postData: map[string]string{
				"foo": "bar",
				"baz": "qux",
			},
			request: &Request{
				Method:   http.MethodPost,
				Resource: "things",
				Protocol: config.ProtocolHTTP,
				Query: url.Values{
					"baz": []string{"qux"},
					"foo": []string{"bar"},
				},
			},
			expectedUri: "things?baz=qux&foo=bar",
		},
		{
			name: "POST with query HTTPS",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
				},
			},
			postData: map[string]string{
				"foo": "bar",
				"baz": "qux",
			},
			request: &Request{
				Method:   http.MethodPost,
				Resource: "things",
				Protocol: config.ProtocolHTTPS,
				Query: url.Values{
					"foo": []string{"bar"},
					"baz": []string{"qux"},
				},
			},
			expectedUri: "things?baz=qux&foo=bar",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Add a body for POST requests.
			if tc.request.Method == http.MethodPost && len(tc.postData) > 0 {
				postBody := new(bytes.Buffer)
				json.NewEncoder(postBody).Encode(tc.postData)

				tc.request.Body = postBody
			}

			err := tc.service.Dial(tc.request)
			if err != nil {
				t.Fatalf("TestCloudService_Dial: %s: %s", tc.name, err)
			}

			if tc.service.CurrentRequest.Method != tc.request.Method {
				t.Errorf("TestService_Dial: %s: expected %v got %v", tc.name, tc.request.Method, tc.service.CurrentRequest.Method)
			}

			expectedUrl := fmt.Sprintf("%s/%s/%s/%s", ts.URL, tc.service.Namespace, tc.service.Name, tc.expectedUri)
			if tc.service.CurrentRequest.URL.String() != expectedUrl {
				t.Errorf("TestService_Dial: %s: expected %v got %v", tc.name, expectedUrl, tc.service.CurrentRequest.URL.String())
			}
		})
	}
}

func ExampleCloudService_Dial() {
	// Instantiate the service.
	myService := &CloudService{
		Service: Service{
			Branch:      "master",
			Environment: "staging",
			Namespace:   "services",
			Name:        "myservice",
		},
		Credentials: &AuthCredentials{
			Email:    "test@test.com",
			Password: "1234",
		},
	}

	// Prepare the request.
	myServiceThingsRequest := &Request{
		Method:   http.MethodGet,
		Resource: "things",
	}

	// Dial it up!
	err := myService.Dial(myServiceThingsRequest)
	if err != nil {
		fmt.Printf("dial err: %s", err)
	}
}

func ExampleCloudService_Call() {
	// Instantiate the service.
	myService := &CloudService{
		Service: Service{
			Branch:      "master",
			Environment: "staging",
			Namespace:   "services",
			Name:        "myservice",
		},
		Credentials: &AuthCredentials{
			Email:    "test@test.com",
			Password: "1234",
		},
	}

	// Prepare the request.
	myServiceThingsRequest := &Request{
		Method:   http.MethodGet,
		Resource: "things",
	}

	// Dial it up!
	err := myService.Dial(myServiceThingsRequest)
	if err != nil {
		fmt.Printf("dial err: %s", err)
	}

	// Do the request.
	myServiceResp, err := myService.Call()
	if err != nil {
		fmt.Printf("call err: %s", err)
	}

	// Make sure we close the body once we're done.
	defer myServiceResp.Body.Close()

	// Decode response.
	serviceResponse := response.Response{}
	jsonErr := json.NewDecoder(myServiceResp.Body).Decode(&serviceResponse)
	if jsonErr != nil {
		fmt.Printf("decode err: %s", err)
	}

	// Handle any error codes.
	switch serviceResponse.Code {
	// Custom error for grants not found.
	case http.StatusNotFound:
		fmt.Println("response err: not found")

		// 200 and 304 are all good.
	case http.StatusOK, http.StatusNotModified:
		break

		// Something somewhere broken!
	default:
		fmt.Println("response err: internal server error")
	}
}
