#!/usr/bin/env bash
#
# Uploads a file to the ginput s3 bucket

if [[ "$#" -ne 1 ]] ; then
  cat << EOF

Uploads a file to the ginput s3 bucket

Usage: $0 <path_to_image>
EOF
  exit 1
fi

bucket=ginput
filepath=$1
filename=$(basename "$1")

aws s3 cp --acl public-read "${filepath}" "s3://${bucket}/${filename}"
