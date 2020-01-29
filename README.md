[![Build Status](https://travis-ci.org/keycloak/keycloak-operator.svg?branch=master)](https://travis-ci.org/keycloak/keycloak-operator)
[![Go Report Card](https://goreportcard.com/badge/github.com/keycloak/keycloak-operator)](https://goreportcard.com/report/github.com/keycloak/keycloak-operator)
[![Coverage Status](https://coveralls.io/repos/github/keycloak/keycloak-operator/badge.svg?branch=master)](https://coveralls.io/github/keycloak/keycloak-operator?branch=master)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# Keycloak Operator
A Kubernetes Operator based on the Operator SDK for creating and syncing resources in Keycloak.

## Help and Documentation

The documentation might be found in the  [docs](./docs/README.asciidoc) directory.

* [Keycloak documentation](https://www.keycloak.org/documentation.html)
* [User Mailing List](https://lists.jboss.org/mailman/listinfo/keycloak-user) - Mailing list for help and general questions about Keycloak
* [JIRA](https://issues.jboss.org/projects/KEYCLOAK) - Issue tracker for bugs and feature requests

## Reporting Security Vulnerabilities

If you've found a security vulnerability, please look at the [instructions on how to properly report it](http://www.keycloak.org/security.html)

## Reporting an issue

If you believe you have discovered a defect in the Keycloak Operator please open an issue in our [Issue Tracker](https://issues.jboss.org/projects/KEYCLOAK).
Please remember to provide a good summary, description as well as steps to reproduce the issue.

## Supported Custom Resources
| *CustomResourceDefinition*                                            | *Description*                                            |
| --------------------------------------------------------------------- | -------------------------------------------------------- |
| [Keycloak](./deploy/crds/keycloak.org_keycloaks_crd.yaml)             | Manages, installs and configures Keycloak on the cluster |
| [KeycloakRealm](./deploy/crds/keycloak.org_keycloakrealms_crd.yaml)   | Represents a realm in a keycloak server                  |
| [KeycloakClient](./deploy/crds/keycloak.org_keycloakclients_crd.yaml) | Represents a client in a keycloak server                 |
| [KeycloakBackup](./deploy/crds/keycloak.org_keycloakbackups_crd.yaml) | Manage Keycloak database backups                         |

## Deploying to a Cluster
*Note*: You will need a running Kubernetes or OpenShift cluster to use the Operator

1. Run `make cluster/prepare` # This will apply the necessary Custom Resource Definitions (CRDs) and RBAC rules to the clusters
2. Run `kubectl apply -f deploy/operator.yaml` # This will start the operator in the current namespace

### Creating Keycloak Instance
Once the CRDs and RBAC rules are applied and the operator is running. Use the examples from the operator.

1. Run `kubectl apply -f deploy/examples/keycloak/keycloak.yaml` 

## Building from Source

To build from source refer to the [building and working with the code base](docs/building.md) guide.

## Components versions

Keycloak Operator supports the following version of key components:

 | *Component*  | *Version/Tag*          |
 | ------------ | ---------------------- |
 | `Keycloak`   | `jboss/keycloak:7.0.1` |
 | `Postgresql` | `9.5`                  |

## Contributing

Before contributing to Keycloak Operator please read our [contributing guidelines](CONTRIBUTING.md).

## Other Keycloak Projects

* [Keycloak](https://github.com/keycloak/keycloak) - Keycloak Server and Java adapters
* [Keycloak Documentation](https://github.com/keycloak/keycloak-documentation) - Documentation for Keycloak
* [Keycloak QuickStarts](https://github.com/keycloak/keycloak-quickstarts) - QuickStarts for getting started with Keycloak
* [Keycloak Docker](https://github.com/jboss-dockerfiles/keycloak) - Docker images for Keycloak
* [Keycloak Gatekeeper](https://github.com/keycloak/keycloak-gatekeeper) - Proxy service to secure apps and services with Keycloak
* [Keycloak Node.js Connect](https://github.com/keycloak/keycloak-nodejs-connect) - Node.js adapter for Keycloak
* [Keycloak Node.js Admin Client](https://github.com/keycloak/keycloak-nodejs-admin-client) - Node.js library for Keycloak Admin REST API

## License

* [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)
