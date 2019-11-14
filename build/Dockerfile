FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV GOCACHE=/go/path/.cache \
    GOPATH=/go/path \
    GOROOT=/usr/local/go
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ENV OPERATOR=/usr/local/bin/keycloak-operator \
    USER_UID=1001 \
    USER_NAME=keycloak-operator

# Because of https://github.com/containers/buildah/issues/1906
# We need to increment this variable whenever we need a rebuild.
# The fix should be implemented very soon
ENV CACHE_BUST=20191114_215444Z

COPY tools /opt/jboss/tools
RUN /opt/jboss/tools/build-operator.sh

COPY bin /usr/local/bin
RUN  /usr/local/bin/user_setup

RUN rm -rf /usr/local/go && rm -rf /go && rm -rf /opt/jboss/tools 

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
