#!/bin/bash -ex

LABELS=$(cat <<-END
LABEL \
com.redhat.component="redhat-sso-7-sso74-rhel8-tech-preview-operator-container"  \
description="Red Hat Single Sign-On 7.4 Operator container image, based on the Red Hat Universal Base Image 8 Minimal container image" \
summary="Red Hat Single Sign-On 7.4 Operator container image, based on the Red Hat Universal Base Image 8 Minimal container image" \
version="7.4" \
io.k8s.description="Operator for Red Hat SSO" \
io.k8s.display-name="Red Hat SSO 7.4 Operator" \
io.openshift.tags="sso,sso74,keycloak,operator" \
name="rh-sso-7-tech-preview\/sso74-rhel8-operator" \
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

if [[ "$BUILD_OPENJ9" != "true" ]]
then
    # remove s390x arch from container.yaml
    # upstream repo should always contain all archs
    sed -i -e '/\-\ s390x/ s/^#*/#/' container.yaml
fi