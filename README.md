# Lush Digital - Microservice Transport (Golang)
A set of convenience structs and interfaces for simplifying transport of data between microservices

## Description
The purpose of the package is to provide a reliable, testable and easy to use means of communicating with microservices
within a service oriented architecture.

## Package Contents
* Transport interface
* Local service struct
* Cloud service struct
* Request struct

## Installation
Install the package as normal:

```bash
$ go get -u github.com/LUSHDigital/microservice-transport-golang
```

## Configuration
There are a few environment variables that can be used to configure this package.

| Variable        | Description                                                                                      |
|-----------------|--------------------------------------------------------------------------------------------------|
| SOA_DOMAIN      | Top level domain of the service environment. Used to build the API gateway URL.                  |
| SOA_GATEWAY_URI | URI of the API gateway e.g. api-gateway                                                          |
| SOA_GATEWAY_URL | Full URL (uri + domain) of the API gateway. Overrides `SOA_DOMAIN` and `SOA_GATEWAY_URI` if set. |

## Documentation
* [General](https://godoc.org/github.com/LUSHDigital/microservice-transport-golang)
* [Config](https://godoc.org/github.com/LUSHDigital/microservice-transport-golang/config)
* [Domain](https://godoc.org/github.com/LUSHDigital/microservice-transport-golang/domain)
* [Errors](https://godoc.org/github.com/LUSHDigital/microservice-transport-golang/errors)
* [Models](https://godoc.org/github.com/LUSHDigital/microservice-transport-golang/models)