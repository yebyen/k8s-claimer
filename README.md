# k8s-claimer

[![Build Status](https://travis-ci.org/deis/k8s-claimer.svg?branch=master)](https://travis-ci.org/deis/k8s-claimer)
[![Go Report Card](https://goreportcard.com/badge/github.com/teamhephy/k8s-claimer)](https://goreportcard.com/report/github.com/teamhephy/k8s-claimer)
[![Docker Repository on Quay](https://quay.io/repository/hephyci/k8s-claimer/status "Docker Repository on Quay")](https://quay.io/repository/hephyci/k8s-claimer)

CLI Downloads:

- [64 Bit Linux](https://storage.googleapis.com/k8s-claimer/k8s-claimer-latest-linux-amd64)
- [32 Bit Linux](https://storage.googleapis.com/k8s-claimer/k8s-claimer-latest-linux-386)
- [64 Bit Mac OS X](https://storage.googleapis.com/k8s-claimer/k8s-claimer-latest-darwin-amd64)
- [32 Bit Mac OS X](https://storage.googleapis.com/k8s-claimer/k8s-claimer-latest-darwin-386)

`k8s-claimer` is a leasing server for a pool of Kubernetes clusters. It will be used as part of our
[deis-workflow end-to-end test](https://github.com/teamhephy/workflow-e2e) infrastructure.

Note that this repository is a work in progress. The code herein is under heavy development,
provides no guarantees and should not be expected to work in any capacity.

As such, it currently does not follow the
[Deis contributing standards](http://docs.deis.io/en/latest/contributing/standards/).

# Design

This server is responsible for holding and managing a set of [Google Container Engine](https://cloud.google.com/container-engine/)
(GKE) clusters. Each cluster can be in the `leased` or `free` state, and this server is responsible for
responding to requests to change a cluster's state, and then safely making the change.

A client who holds the lease for a cluster has a [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier)
indicating their ownership as well as the guarantee that nobody else will get the lease before
either their lease duration expires or someone releases the lease with their UUID. The client
specifies the lease duration when they acquire it.

For implementation details, see [the architecture document](doc/architecture.md)

## Azure:
* You need to create an SP that has access to the API. You can do that by issuing the following command: 
```
az ad sp create-for-rbac --role="Contributor" --scopes="/subscriptions/<SUBSCRIPTION_ID>" --name="<NAME OF SP>"
```
* We currently need to SCP the kubeconfig file off of the master node so you need to provide k8s-claimer with the ssh key (private) used to setup your leasable clusters (this means they all need to be the same).
* You cannot lease a cluster by version since the Azure API does not return the version of the API running.


# Configuration
This server is configured exclusively with environment variables. See below for a list and
description of each.

| ENV Var | Description |
|---------|-------------|
| BIND_PORT | The port to bind the server to. Defaults to `8080` |
| BIND_HOST | The host to bind the server to. Defaults to `0.0.0.0`
| NAMESPACE | The namespace in which to store lease data (lease data is stored on annotations on a service in this namespace). Defaults to `k8s-claimer` | 
| SERVICE_NAME | The service on which to store lease data. Defaults to `k8s-claimer` |
| AUTH_TOKEN | The authentication token that clients must use to acquire and release leases |
| GOOGLE_CLOUD_ACCOUNT_FILE | Base64 encoded JWT of the Service Account for GKE | 
| GOOGLE_CLOUD_PROJECT_ID | The Google Cloud project ID for the project that holds the GKE clusters to lease. This is a required field |
| GOOGLE_CLOUD_ZONE | The zone that clusters can be leased from. Pass `-` to indicate all zones. Defaults to `-` | 
| AZURE_CLIENT_ID | The Service Principal ID used to make Azure API calls |
| AZURE_CLIENT_SECRET | The secret for the Service Principal | 
| AZURE_TENANT_ID | The tenant of the Service Principal | 
| AZURE_SUBSCRIPTION_ID | The subscription where the leasable clusters live |


## GOOGLE_CLOUD_ACCOUNT_FILE
You can get a JWT file from the Google Cloud Platform console by following these steps:
  - Go to `Permissions`
  - Create a service account if necessary
  - Click on the three vertical dots to the far right of the service account row for which you would like permissions
  - Click `Create key`
  - Make sure `JSON` is checked and click `Create`

# CLI
A command line interface is provided which talks to the REST API server. You can download the binaries
using the links above, or build it yourself by executing `make bootstrap build-cli`.

See `k8s-claimer --help` for full CLI documentation.

```shell
$ k8s-claimer --help
NAME:
   k8s-claimer - This CLI can be used against a k8s-claimer server to acquire and release leases

USAGE:
   k8s-claimer [global options] command [command options] [arguments...]

COMMANDS:
     lease
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --server value  The k8s-claimer server to talk to
   --help, -h      show help
   --version, -v   print the version
```

## Create a Lease

```shell
$ k8s-claimer lease create --help
NAME:
   k8s-claimer lease create - Creates a new lease and returns 'export' statements to set the lease values as environment variables. Set the 'env-prefix' flag to prefix the environment variable names. If you pass that flag, a '_' character will separate the prefix with the rest of the environment variable name. Below are the basic environment variable names:

- IP - the IP address of the Kubernetes master server
- TOKEN - contains the lease token. Use this when you run 'k8s-claimer-cli lease delete'
- CLUSTER_NAME - contains the name of the cluster. For informational purposes only

The Kubeconfig file will be written to kubeconfig-file


USAGE:
   k8s-claimer lease create [command options] [arguments...]

OPTIONS:
   --duration value         The duration of the lease in seconds (default: 10)
   --env-prefix value       The prefix for all environment variables that this command sets
   --kubeconfig-file value  The location of the resulting Kubeconfig file (default: "./kubeconfig.yaml")
   --cluster-regex value    A regular expression that will be used to match which cluster you lease
   --cluster-version value  A version string that will be used to find a cluster to lease
   --provider value         Which cloud provider to use when creating a cluster lease. Acceptable values are azure and google. If a value is not provided it will return an error.
```

Example
```shell
$ k8s-claimer --server <server-name> lease create
export IP="1.2.3.4"
export TOKEN="<token>"
export CLUSTER_NAME="cattier"
```

## Delete a Lease

```shell
$ k8s-claimer lease delete --help
NAME:
   k8s-claimer lease delete - Releases a currently held lease. Pass the lease token as the first and only parameter to this command. For example:

k8s-claimer-cli lease delete $TOKEN


USAGE:
   k8s-claimer lease delete [command options] [arguments...]

OPTIONS:
   --provider value  Which cloud provider to use when deleting a cluster lease. Acceptable values are azure and google. If a value is not provided it will return an error.
```

Example
```shell
$ k8s-claimer --server <server-name> --provider GKE lease delete <token>
Deleted lease <token>
```

# API

The server exposes a REST API to acquire and release leases for clusters. The subsections
herein list each endpoint.

## `POST /lease`

Acquire a new lease.

### Request Body

```json
{"max_time": 30}
```

Note that the value of `max_time` is the maximum lease duration in seconds. It must be a number.
After this duration expires, the lease will be automatically released.

### Responses

Unless otherwise noted, all responses except for `200 OK` indicate that the lease was not acquired.

In non-200 response code scenarios, a body may be returned with an explanation of the error,
but the existence or contents of that body are not guaranteed.

#### `401 Bad Request`

This response code is returned with no specific body if the request body was malformed.

#### `500 Internal Server Error`

This response code is returned if any of the following occur:

- The server couldn't communicate with the Kubernetes Master to get the service object
- The server couldn't communicate with the GKE API
- A cluster was available, but the new lease information couldn't be saved
- An expired lease exists but it points to a non-existent cluster
- The lease was succesful but the response body couldn't be rendered

#### `409 Conflict`

This response code is returned if there are no clusters available for lease.

#### `200 OK`

This response code is returned along with the below response body if a lease was successfully
acquired.

```json
{
  "kubeconfig": "RFC 4648 base64 encoded Kubernetes config file. After decoding, this value can be written to ~/.kube/config for use with kubectl",
  "ip": "The IP address of the Kubernetes master server in GKE",
  "token": "The token of the lease. This is your proof of ownership of the cluster, until the lease expires or you release it",
  "cluster_name": "The name of the cluster. This value is purely informational, and fetched from GKE"
}
```

## `Delete /lease/{token}`

Release an existing lease, identified by `{token}`.

### Responses

All responses except for `200 OK` indicate that no leases were changed. Since the state of the
lease identified by `{token}` (if there was one) can change over time, there are no guarantees
on the state of the lease after this API call returns in these cases.

In all cases, a body may be returned with an explanation of the response, but the existece or
contents of that body are not guaranteed.

#### `401 Bad Request`

This response code is returned in the following cases:

- The URL path did not include a lease token
- The lease token was malformed

#### `500 Internal Server Error`

This response code is returned in the following cases:

- The server couldn't communicate with the Kubernetes Master to get the service object
- The lease was found and deleted, but the updated lease statuses couldn't be saved

#### `409 Conflict`

This response code is returned when no lease exists with the given token.

#### `200 OK`

The lease was successfully deleted. The given token is no longer valid and should not be reused.
