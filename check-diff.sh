#!/bin/bash -ex

if [ -n "$(git status --untracked-files=no --porcelain)" ]; then
  echo "code needs to be regenerated: use \"make setup\" and check in any resulting diffs"
  git --no-pager diff
  exit 1
fi

