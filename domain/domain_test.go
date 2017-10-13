package domain

import (
	"fmt"
	"testing"
)

func TestBuildServiceDNSName(t *testing.T) {
	tt := []struct {
		name             string
		service          string
		branch           string
		environment      string
		serviceNamespace string
		expectedDnsName  string
	}{
		{
			name:             "Normal data",
			service:          "test",
			branch:           "master",
			environment:      "staging",
			serviceNamespace: "test",
			expectedDnsName:  "test-master-staging.test",
		},
		{
			name:             "Extreme data",
			service:          "21323kl1j3913issvxc9vx0",
			branch:           "(!()*)(*!KJ",
			environment:      "sljsjlfjdkgj",
			serviceNamespace: ")ID`hdfy7d7f",
			expectedDnsName:  "21323kl1j3913issvxc9vx0-(!()*)(*!KJ-sljsjlfjdkgj.)ID`hdfy7d7f",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actualDNSName := BuildServiceDNSName(tc.service, tc.branch, tc.environment, tc.serviceNamespace)
			if actualDNSName != tc.expectedDnsName {
				t.Errorf("TestBuildServiceDNSName: %s: expected %v got %v", tc.name, tc.expectedDnsName, actualDNSName)
			}
		})
	}
}

func ExampleBuildServiceDNSName() {
	dnsName := BuildServiceDNSName("myservice", "master", "staging", "services")
	fmt.Println(dnsName)

	// Output: myservice-master-staging.services
}

func TestBuildCloudServiceUrl(t *testing.T) {
	tt := []struct {
		name                    string
		apiGatewayUrl           string
		serviceNamespace        string
		serviceName             string
		expectedCloudServiceUrl string
	}{
		{
			name:                    "Normal data",
			apiGatewayUrl:           "test.com",
			serviceNamespace:        "test",
			serviceName:             "test",
			expectedCloudServiceUrl: "test.com/test/test",
		},
		{
			name:                    "Extreme data",
			apiGatewayUrl:           "te(SDS(sdsdsdst.com",
			serviceNamespace:        "sdfisfpsif9((DF",
			serviceName:             "D&D&*FDHFHSDFHDF",
			expectedCloudServiceUrl: "te(SDS(sdsdsdst.com/sdfisfpsif9((DF/D&D&*FDHFHSDFHDF",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actualCloudServiceUrl := BuildCloudServiceUrl(tc.apiGatewayUrl, tc.serviceNamespace, tc.serviceName)
			if actualCloudServiceUrl != tc.expectedCloudServiceUrl {
				t.Errorf("TestBuildServiceDNSName: %s: expected %v got %v", tc.name, tc.expectedCloudServiceUrl, actualCloudServiceUrl)
			}
		})
	}
}

func ExampleBuildCloudServiceUrl() {
	cloudServiceUrl := BuildCloudServiceUrl("my-api-gateway.com", "services", "myservice")
	fmt.Println(cloudServiceUrl)

	// Output: my-api-gateway.com/services/myservice
}
