# Building from source

## Local Development
*Note*: You will need a running Kubernetes or OpenShift cluster to use the Operator

1.  clone this repo to `$GOPATH/src/github.com/keycloak/keycloak-operator`
2.  run `make setup/mod cluster/prepare`
3.  run `make code/run`
-- The above step will launch the operator on the local machine
-- To see how do debug the operator or how to deploy to a cluster, see below alternatives to step 3
4. In a new terminal run `make cluster/create/examples`

To clean the cluster (Removes CRDs, CRs, RBAC and namespace)
1. run `make cluster/clean`

### Alternative Step 2: Debug in Goland
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

### Alternative Step 3: Debug in VSCode
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

### Alternative Step 3: Deploying to a Cluster
Deploy the operator into the running cluster
1. build image with `operator-sdk build <image registry>/<organisation>/keycloak-operator:<tag>`. e.g. `operator-sdk build quay.io/keycloak/keycloak-operator:test`
2. Change the `image` property in `deploy/operator.yaml` to the above full image path
3. run `kubectl apply -f deploy/operator.yaml -n <NAMESPACE>`


### Makefile command reference
#### Operator Setup Management
| *Command*                      | *Description*                                                                                          |
| ------------------------------ | ------------------------------------------------------------------------------------------------------ |
| `make cluster/prepare`         | Creates the `keycloak` namespace, applies all CRDs to the cluster and sets up the RBAC files           |
| `make cluster/clean`           | Deletes the `keycloak` namespace, all `keycloak.org` CRDs and all RBAC files named `keycloak-operator` |
| `make cluster/create/examples` | Applies the example Keycloak and KeycloakRealm CRs                                                     |

#### Tests
| *Command*                    | *Description*                                           |
| ---------------------------- | ------------------------------------------------------- |
| `make test/unit`             | Runs unit tests                                         |
| `make test/e2e`              | Runs e2e tests                                          |
| `make test/coverage/prepare` | Prepares coverage report from unit and e2e test results |
| `make test/coverage`         | Generates coverage report                               |

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

#### CI
| *Command*           | *Description*                                                              |
| ------------------- | -------------------------------------------------------------------------- |
| `make setup/travis` | Downloads operator-sdk, makes it executable and copys to `/usr/local/bin/` |

