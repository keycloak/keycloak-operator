FROM registry.ci.openshift.org/openshift/release:golang-1.13 AS build-env

COPY . /src/

RUN cd /src && \
    make code/compile && \
    echo "Build SHA1: $(git rev-parse HEAD)" && \
    echo "$(git rev-parse HEAD)" > /src/BUILD_INFO

# final stage
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

##LABELS

RUN microdnf update && microdnf clean all && rm -rf /var/cache/yum/*

COPY --from=build-env /src/BUILD_INFO /src/BUILD_INFO
COPY --from=build-env /src/tmp/_output/bin/keycloak-operator /usr/local/bin

ENTRYPOINT ["/usr/local/bin/keycloak-operator"]
