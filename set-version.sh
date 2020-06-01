#!/bin/bash -ex

# Implemented to deal with patterns such Version = "x.y.z"
# An additional advantage of this approach is that it preserves formatting
# (which is desirable when working with Go)
function replace_value_in_file() {
  local string_to_be_replaced=$1
  local replacement=$2
  local file=$3
  local quote=${4:-}

  sed -i "s/\($string_to_be_replaced\)\(.*\).*/\1${quote}${replacement}${quote}/" $file
}

replace_value_in_file ".*Version = " "$1" "version/version.go" "\""
replace_value_in_file ".*DefaultKeycloakImage.*= " "quay.io\/keycloak\/keycloak:$1" "pkg/model/image_manager.go" "\""
replace_value_in_file ".*image: " "quay.io\/keycloak\/keycloak-operator:$1" "deploy/operator.yaml"
