package microservicetransport

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"fmt"

	"github.com/LUSHDigital/microservice-core-golang/response"
	"github.com/LUSHDigital/microservice-transport-golang/config"
)

func TestService_Dial(t *testing.T) {
	tt := []struct {
		name           string
		service        Service
		request        *Request
		postData       map[string]string
		expectedMethod string
		expectedUrl    string
		expectedBody   string
	}{
		{
			name: "GET HTTP",
			service: Service{
				Branch:      "master",
				Environment: "staging",
				Namespace:   "services",
				Name:        "myservice",
			},
			request: &Request{
				Method:   http.MethodGet,
				Resource: "things",
				Protocol: config.ProtocolHTTP,
			},
			expectedUrl: "http://myservice-master-staging.myservice/things",
		},
		{
			name: "GET HTTPS",
			service: Service{
				Branch:      "master",
				Environment: "staging",
				Namespace:   "services",
				Name:        "myservice",
			},
			request: &Request{
				Method:   http.MethodGet,
				Resource: "things",
				Protocol: config.ProtocolHTTPS,
			},
			expectedUrl: "https://myservice-master-staging.myservice/things",
		},
		{
			name: "GET with query HTTP",
			service: Service{
				Branch:      "master",
				Environment: "staging",
				Namespace:   "services",
				Name:        "myservice",
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
			expectedUrl: "http://myservice-master-staging.myservice/things?baz=qux&foo=bar",
		},
		{
			name: "GET with query HTTPS",
			service: Service{
				Branch:      "master",
				Environment: "staging",
				Namespace:   "services",
				Name:        "myservice",
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
			expectedUrl: "https://myservice-master-staging.myservice/things?baz=qux&foo=bar",
		},
		{
			name: "POST HTTP",
			service: Service{
				Branch:      "master",
				Environment: "staging",
				Namespace:   "services",
				Name:        "myservice",
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
			expectedUrl: "http://myservice-master-staging.myservice/things",
		},
		{
			name: "POST HTTPS",
			service: Service{
				Branch:      "master",
				Environment: "staging",
				Namespace:   "services",
				Name:        "myservice",
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
			expectedUrl: "https://myservice-master-staging.myservice/things",
		},
		{
			name: "POST with query HTTP",
			service: Service{
				Branch:      "master",
				Environment: "staging",
				Namespace:   "services",
				Name:        "myservice",
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
			expectedUrl: "http://myservice-master-staging.myservice/things?baz=qux&foo=bar",
		},
		{
			name: "POST with query HTTPS",
			service: Service{
				Branch:      "master",
				Environment: "staging",
				Namespace:   "services",
				Name:        "myservice",
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
			expectedUrl: "https://myservice-master-staging.myservice/things?baz=qux&foo=bar",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Add a body for POST requests.
			if tc.request.Method == http.MethodPost && len(tc.postData) > 0 {
				postBody := new(bytes.Buffer)
				json.NewEncoder(postBody).Encode(tc.postData)

				tc.request.Body = ioutil.NopCloser(postBody)
			}

			err := tc.service.Dial(tc.request)
			if err != nil {
				t.Errorf("TestService_Dial: %s: %s", tc.name, err)
			}

			if tc.service.CurrentRequest.Method != tc.request.Method {
				t.Errorf("TestService_Dial: %s: expected %v got %v", tc.name, tc.request.Method, tc.service.CurrentRequest.Method)
			}

			if tc.service.CurrentRequest.URL.String() != tc.expectedUrl {
				t.Errorf("TestService_Dial: %s: expected %v got %v", tc.name, tc.expectedUrl, tc.service.CurrentRequest.URL.String())
			}
		})
	}
}

func TestService_GetName(t *testing.T) {
	tt := []struct {
		name           string
		service        Service
		expectedName string
	}{
		{
			name: "Normal",
			service: Service{
				Branch:      "master",
				Environment: "staging",
				Namespace:   "services",
				Name:        "myservice",
			},
			expectedName: "myservice",
		},
		{
			name: "Crazy",
			service: Service{
				Branch:      "massdsdfsdjf89uter",
				Environment: "sdfsdf34341",
				Namespace:   "l1j2312klj3k21j3",
				Name:        "-sf9s9f9ds0f9-",
			},
			expectedName: "-sf9s9f9ds0f9-",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.service.GetName() != tc.expectedName {
				t.Errorf("TestService_GetName: %s: expected %v got %v", tc.name, tc.expectedName, tc.service.GetName())
			}
		})
	}
}

func ExampleService_Dial() {
	// Instantiate the service.
	myService := &Service{
		Branch:      "master",
		Name:        "myservice",
		Environment: "staging",
		Namespace:   "services",
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

func ExampleService_Call() {
	// Instantiate the service.
	myService := &Service{
		Branch:      "master",
		Name:        "myservice",
		Environment: "staging",
		Namespace:   "services",
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
