# This Dockerfile is used for building Red Hat Single Sign-On Operator.
# If you're looking for a Dockerfile for Keycloak, please look at Dockerfile.keycloak

FROM ubi8-minimal:8-released AS build-env

RUN microdnf install -y git make golang

COPY . /src/
RUN cd /src && make code/compile
RUN cd /src && echo "Build SHA1: $(git rev-parse HEAD)"
RUN cd /src && echo "$(git rev-parse HEAD)" > /src/BUILD_INFO

# final stage
FROM ubi8-minimal:8-released

LABEL \
    com.redhat.component="redhat-sso-7-sso74-operator-rhel8-container"  \
    description="Red Hat Single Sign-On 7.4 Operator on OpenJDK OpenShift container image, based on the Red Hat Universal Base Image 8 Minimal container image" \
    summary="Red Hat Single Sign-On 7.4 Operator on OpenJDK OpenShift container image, based on the Red Hat Universal Base Image 8 Minimal container image" \
    version="7.4" \
    io.k8s.description="Operator for Red Hat SSO" \
    io.k8s.display-name="Red Hat SSO 7.4 Operator" \
    io.openshift.tags="sso,sso74,keycloak,operator" \
    name="rh-sso-7/sso74-operator-rhel8" \
    maintainer="Red Hat Single Sign-On Team"

COPY --from=build-env /src/BUILD_INFO /src/BUILD_INFO
COPY --from=build-env /src/tmp/_output/bin/keycloak-operator /usr/local/bin

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
