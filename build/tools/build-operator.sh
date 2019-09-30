#!/bin/sh

# Default Keycloak repository and destination path
# In the future we may want to provide it as arguments to the build

GIT_REPO="https://github.com/keycloak/keycloak-operator.git"
KEYCLOAK_PATH="/go/src/github.com/keycloak"

mkdir -p $KEYCLOAK_PATH
microdnf update && microdnf install -y tar gzip git make && microdnf clean all && rm -rf /var/cache/yum/*

# Install Go
VERSION="1.12.1"
curl https://storage.googleapis.com/golang/go$VERSION.linux-amd64.tar.gz | tar -C /usr/local -xzf - \
  && rm -rf $GOROOT/{pkg/linux_amd64_race,test,doc,api}/* \
  $GOROOT/pkg/tool/linux_amd64/{vet,doc,cover,trace,nm,fix,test2json,objdump} \
  $GOROOT/bin/godoc \
  /var/lib/rpm/Packages && \
  find /usr/share/locale/ -name tar.mo | xargs rm && \
  mkdir -p $GOPATH/bin && \
  chmod g+xw -R /go && \
  chmod g+xw -R $(go env GOROOT)

# Clone Keycloak repository
git clone --depth 1 $GIT_REPO $KEYCLOAK_PATH/keycloak-operator

# Build and copy the binary
cd /go/src/github.com/keycloak/keycloak-operator && make code/compile
cp ./tmp/_output/bin/keycloak-operator /usr/local/bin
