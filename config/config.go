package config

const (
	// AuthConsumerKey - Name of the context key used to pass consumer information
	// between layers of middleware.
	AuthConsumerKey = "auth_consumer"

	// AuthHeader - Name of the HTTP header to use for authentication and
	// authorization.
	AuthHeader = "Authorization"

	// AuthHeaderPrefix - Prefix expected for the HTTP auth header value.
	AuthHeaderPrefix = "Bearer "

	// ContentTypeHeader - Name of the HTTP header to use for content types.
	ContentTypeHeader = "Content-Type"

	// ProtocolHTTP - Protocol string for non-ssl requests.
	ProtocolHTTP = "http"

	// ProtocolHTTPS - Protocol string for ssl requests.
	ProtocolHTTPS = "https"

	// ResourceGrantsKey - Name of the context key used to pass resource grants
	// between layers of middleware.
	ResourceGrantsKey = "resource_grants"

	// RequestKey - Name of the context key used to pass a service request
	// between layers of middleware.
	RequestKey = "request"

	// ServiceBranchHeader - Name of the HTTP header to use for service branch.
	ServiceBranchHeader = "x-service-branch"

	// ServiceEnvironmentHeader - Name of the HTTP header to use for service environment.
	ServiceEnvironmentHeader = "x-service-environment"

	// ServiceVersionHeader - Name of the HTTP header to use for service version.
	ServiceVersionHeader = "x-service-version"

	// ServiceKey - Name of the context key used to pass service details
	// between layers of middleware.
	ServiceKey = "service"

	// AggregatorDomainPrefix - The prefix value used for aggregator domains.
	AggregatorDomainPrefix = "agg"
)

var (
	// ServiceBranch - The default VCS branch to use for services.
	ServiceBranch = "master"
)
