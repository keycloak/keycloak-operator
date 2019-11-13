FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV GOCACHE=/go/path/.cache \
    GOPATH=/go/path \
    GOROOT=/usr/local/go
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ENV OPERATOR=/usr/local/bin/keycloak-operator \
    USER_UID=1001 \
    USER_NAME=keycloak-operator

ENV SHA1_FILE=/BUILD_SHA1.txt
ENV BUILD_DATE_FILE=/BUILD_DATE.txt

RUN echo "$(date)" > $BUILD_DATE_FILE

COPY tools /opt/jboss/tools
RUN /opt/jboss/tools/build-operator.sh

COPY bin /usr/local/bin
RUN  /usr/local/bin/user_setup

RUN rm -rf /usr/local/go && rm -rf /go && rm -rf /opt/jboss/tools 

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
