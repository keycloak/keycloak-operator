#!/bin/bash -ex

LABELS=$(cat <<-END
LABEL \
com.redhat.component="redhat-sso-7-rhel8-operator-container"  \
description="Red Hat Single Sign-On 7. Operator container image, based on the Red Hat Universal Base Image 8 Minimal container image" \
summary="Red Hat Single Sign-On 7.5 Operator container image, based on the Red Hat Universal Base Image 8 Minimal container image" \
version="7.5" \
io.k8s.description="Operator for Red Hat SSO" \
io.k8s.display-name="Red Hat SSO 7.5 Operator" \
io.openshift.tags="sso,sso75,keycloak,operator" \
name="rh-sso-7\/sso7-rhel8-operator" \
maintainer="Red Hat Single Sign-On Team"
END
)

sed -i \
    -e 's/registry.ci.openshift.org\/openshift\/release:golang-1.13/openshift\/golang-builder:1.13/' \
    -e 's/FROM registry.access.redhat.com/FROM registry.redhat.io/' \
    -e 's/COPY . /COPY keycloak-operator-*.tar.gz /' \
    -e 's,RUN cd /src ,RUN cd /src \&\& tar -x --strip-components=1 -f keycloak-operator-*.tar.gz ,' \
    -e "s/##LABELS/$LABELS/g" \
    Dockerfile
