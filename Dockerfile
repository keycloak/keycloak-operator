FROM registry.access.redhat.com/ubi8/ubi-minimal:8.1 AS build-env

RUN microdnf install -y git make golang

ADD ./ /src
RUN cd /src && make code/compile
RUN cd /src && echo "Build SHA1: $(git rev-parse HEAD)"
RUN cd /src && echo "$(git rev-parse HEAD)" > /src/BUILD_INFO

# final stage
FROM registry.access.redhat.com/ubi8/ubi-minimal:8.1

ENV OPERATOR=/usr/local/bin/keycloak-operator \
    USER_UID=1001 \
    USER_NAME=keycloak-operator

RUN microdnf update && microdnf clean all && rm -rf /var/cache/yum/*

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

COPY --from=build-env /src/BUILD_INFO /src/BUILD_INFO
COPY --from=build-env /src/tmp/_output/bin/keycloak-operator /usr/local/bin

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
