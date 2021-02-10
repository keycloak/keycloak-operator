[![Build Status](https://travis-ci.org/keycloak/keycloak-operator.svg?branch=master)](https://travis-ci.org/keycloak/keycloak-operator)
[![Go Report Card](https://goreportcard.com/badge/github.com/keycloak/keycloak-operator)](https://goreportcard.com/report/github.com/keycloak/keycloak-operator)
[![Coverage Status](https://coveralls.io/repos/github/keycloak/keycloak-operator/badge.svg?branch=master)](https://coveralls.io/github/keycloak/keycloak-operator?branch=master)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# Keycloak Operator
A Kubernetes Operator based on the Operator SDK for creating and syncing resources in Keycloak.

## Help and Documentation

The official documentation might be found in the [here](https://www.keycloak.org/docs/latest/server_installation/index.html#_operator).

* [Keycloak documentation](https://www.keycloak.org/documentation.html)
* [User Mailing List](https://lists.jboss.org/mailman/listinfo/keycloak-user) - Mailing list for help and general questions about Keycloak
* [JIRA](https://issues.redhat.com/browse/KEYCLOAK-16220?jql=project%20%3D%20KEYCLOAK%20AND%20component%20%3D%20%22Container%20-%20Operator%22%20ORDER%20BY%20updated%20DESC) - Issue tracker for bugs and feature requests

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


## Deployment to a Kubernetes or Openshift cluster

The official documentation contains installation instruction for this Operator.

[Getting started with keycloak-operator on Openshift](https://www.keycloak.org/getting-started/getting-started-operator-openshift)

[Getting started with keycloak-operator on Kubernetes](https://www.keycloak.org/getting-started/getting-started-operator-kubernetes)

[Operator installation](https://www.keycloak.org/docs/latest/server_installation/index.html#_installing-operator)


## Developer Reference
*Note*: You will need a running Kubernetes or OpenShift cluster to use the Operator

1. Run `make cluster/prepare` # This will apply the necessary Custom Resource Definitions (CRDs) and RBAC rules to the clusters
2. Run `kubectl apply -f deploy/operator.yaml` # This will start the operator in the current namespace

### Creating Keycloak Instance
Once the CRDs and RBAC rules are applied and the operator is running. Use the examples from the operator.

1. Run `kubectl apply -f deploy/examples/keycloak/keycloak.yaml`

### Local Development
*Note*: You will need a running Kubernetes or OpenShift cluster to use the Operator

1.  clone this repo to `$GOPATH/src/github.com/keycloak/keycloak-operator`
2.  run `make setup/mod cluster/prepare`
3.  run `make code/run`
-- The above step will launch the operator on the local machine
-- To see how do debug the operator or how to deploy to a cluster, see below alternatives to step 3
4. In a new terminal run `make cluster/create/examples`
5. Optional: configure Ingress and DNS Resolver
   - minikube: \
     -- run `minikube addons enable ingress` \
     -- run `./hack/modify_etc_hosts.sh`
   - Docker for Mac: \
     -- run `kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-0.32.0/deploy/static/provider/cloud/deploy.yaml`
        (see also https://kubernetes.github.io/ingress-nginx/deploy/) \
     -- run `./hack/modify_etc_hosts.sh keycloak.local 127.0.0.1`
6. Run `make test/e2e`

To clean the cluster (Removes CRDs, CRs, RBAC and namespace)
1. run `make cluster/clean`

#### Alternative Step 2: Debug in Goland
Debug the operator in [Goland](https://www.jetbrains.com/go/)
1. go get -u github.com/go-delve/delve/cmd/dlv
2. Create new `Go Build` debug configuration
3. Change the properties to the following
```
* Name = Keycloak Operator
* Run Kind = File
* Files = <project full path>/cmd/manager/main.go
* Working Directory = <project full path>
* Environment = KUBERNETES_CONFIG=<kube config path>;WATCH_NAMESPACE=keycloak
```
3. Apply and click Debug Keycloak operator

#### Alternative Step 3: Debug in VSCode
Debug the operator in [VS Code](https://code.visualstudio.com/docs/languages/go)
1. go get -u github.com/go-delve/delve/cmd/dlv
2. Create new launch configuration, changing your kube config location
```json
{
  "name": "Keycloak Operator",
  "type": "go",
  "request": "launch",
  "mode": "auto",
  "program": "${workspaceFolder}/cmd/manager/main.go",
  "env": {
    "WATCH_NAMESPACE": "keycloak",
    "KUBERNETES_CONFIG": "<kube config path>"
  },
  "cwd": "${workspaceFolder}",
  "args": []
}
```
3. Debug Keycloak Operator

#### Alternative Step 3: Deploying to a Cluster
Deploy the operator into the running cluster
1. build image with `operator-sdk build <image registry>/<organisation>/keycloak-operator:<tag>`. e.g. `operator-sdk build quay.io/keycloak/keycloak-operator:test`
2. Change the `image` property in `deploy/operator.yaml` to the above full image path
3. run `kubectl apply -f deploy/operator.yaml -n <NAMESPACE>`

#### Alternative Step 6: Debug the e2e tests in Goland
Debug the e2e operator tests in [Goland](https://www.jetbrains.com/go/)
1. Set `Test kind` to `Package`
2. Set `Working directory` to `<your project directory>`
3. Set `Go tool arguments` to `-i -parallel=1`
4. Set `Program arguments` to `-root=<your project directory> -kubeconfig=<your home directory>/.kube/config -globalMan deploy/empty-init.yaml -namespacedMan deploy/empty-init.yaml -test.v -singleNamespace -localOperator -test.timeout 0`
5. Apply and click Debug Keycloak operator

### Makefile command reference
#### Operator Setup Management
| *Command*                      | *Description*                                                                                          |
| ------------------------------ | ------------------------------------------------------------------------------------------------------ |
| `make cluster/prepare`         | Creates the `keycloak` namespace, applies all CRDs to the cluster and sets up the RBAC files           |
| `make cluster/clean`           | Deletes the `keycloak` namespace, all `keycloak.org` CRDs and all RBAC files named `keycloak-operator` |
| `make cluster/create/examples` | Applies the example Keycloak and KeycloakRealm CRs                                                     |

#### Tests
| *Command*                    | *Description*                                               |
| ---------------------------- | ----------------------------------------------------------- |
| `make test/unit`             | Runs unit tests                                             |
| `make test/e2e`              | Runs e2e tests with operator ran locally                    |
| `make test/e2e-latest-image` | Runs e2e tests with latest available operator image running in the cluster |
| `make test/e2e-local-image`  | Runs e2e tests with local operator image running in the cluster |
| `make test/coverage/prepare` | Prepares coverage report from unit and e2e test results     |
| `make test/coverage`         | Generates coverage report                                   |

##### Running tests without cluster admin permissions
It's possible to deploy CRDs, roles, role bindings, etc. separately from running the tests:
1. Run `make cluster/prepare` as a cluster admin.
2. Run `make test/ibm-validation` as a user. The user needs the following permissions to run te tests:
```
apiGroups: ["", "apps", "keycloak.org"]
resources: ["persistentvolumeclaims", "deployments", "statefulsets", "keycloaks", "keycloakrealms", "keycloakusers", "keycloakclients", "keycloakbackups"]
verbs: ["*"]
```
Please bear in mind this is intended to be used for internal purposes as there's no guarantee it'll work without any issues.

#### Local Development
| *Command*                 | *Description*                                                                    |
| ------------------------- | -------------------------------------------------------------------------------- |
| `make setup`              | Runs `setup/mod` `setup/githooks` `code/gen`                                     |
| `make setup/githooks`     | Copys githooks from `./githooks` to `.git/hooks`                                 |
| `make setup/mod`          | Resets the main module's vendor directory to include all packages                |
| `make setup/operator-sdk` | Installs the operator-sdk                                                        |
| `make code/run`           | Runs the operator locally for development purposes                               |
| `make code/compile`       | Builds the operator                                                              |
| `make code/gen`           | Generates/Updates the operator files based on the CR status and spec definitions |
| `make code/check`         | Checks for linting errors in the code                                            |
| `make code/fix`           | Formats code using [gofmt](https://golang.org/cmd/gofmt/)                        |
| `make code/lint`          | Checks for linting errors in the code                                            |
| `make client/gen`         | Generates/Updates the clients bases on the CR status and spec definitions        |

#### Application Monitoring

NOTE: This functionality works only in OpenShift environment.

| *Command*                         | *Description*                                           |
| --------------------------------- | ------------------------------------------------------- |
| `make cluster/prepare/monitoring` | Installs and configures Application Monitoring Operator |

#### CI
| *Command*           | *Description*                                                              |
| ------------------- | -------------------------------------------------------------------------- |
| `make setup/travis` | Downloads operator-sdk, makes it executable and copys to `/usr/local/bin/` |

#### Components versions

All images used by the Operator might be controlled using dedicated Environmental Variables:

 | *Image*             | *Environment variable*          | *Default*                                                        |
 | ------------------- | ------------------------------- | ---------------------------------------------------------------- |
 | `Keycloak`          | `RELATED_IMAGE_KEYCLOAK`                | `quay.io/keycloak/keycloak:9.0.2`                                |
 | `RHSSO` for OpenJ9  | `RELATED_IMAGE_RHSSO_OPENJ9`            | `registry.redhat.io/rh-sso-7/sso74-openshift-rhel8:7.4-1`        |
 | `RHSSO` for OpenJDK | `RELATED_IMAGE_RHSSO_OPENJDK`           | `registry.redhat.io/rh-sso-7/sso74-openshift-rhel8:7.4-1`        |
 | Init container      | `RELATED_IMAGE_KEYCLOAK_INIT_CONTAINER` | `quay.io/keycloak/keycloak-init-container:master`                |
 | Backup container    | `RELATED_IMAGE_RHMI_BACKUP_CONTAINER`   | `quay.io/integreatly/backup-container:1.0.16`                    |
 | Postgresql          | `RELATED_IMAGE_POSTGRESQL`              | `registry.redhat.io/rhel8/postgresql-10:1`                       |

## Contributing

Before contributing to Keycloak Operator please read our [contributing guidelines](CONTRIBUTING.md).

## Other Keycloak Projects

* [Keycloak](https://github.com/keycloak/keycloak) - Keycloak Server and Java adapters
* [Keycloak Documentation](https://github.com/keycloak/keycloak-documentation) - Documentation for Keycloak
* [Keycloak QuickStarts](https://github.com/keycloak/keycloak-quickstarts) - QuickStarts for getting started with Keycloak
* [Keycloak Docker](https://github.com/jboss-dockerfiles/keycloak) - Docker images for Keycloak
* [Keycloak Node.js Connect](https://github.com/keycloak/keycloak-nodejs-connect) - Node.js adapter for Keycloak
* [Keycloak Node.js Admin Client](https://github.com/keycloak/keycloak-nodejs-admin-client) - Node.js library for Keycloak Admin REST API

## License

* [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)
