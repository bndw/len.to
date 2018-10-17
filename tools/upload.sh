#!/usr/bin/env bash
#
# Uploads an image to the len.to s3 bucket.
set -e

if [[ "$#" -ne 1 ]] ; then
  cat << EOF

Uploads a file to the len.to s3 bucket.

Usage: $0 <path_to_image>
EOF
  exit 1
fi

bucket=ginput
filepath=$1
filename=$(basename "$1")
url="https://d17enza3bfujl8.cloudfront.net/${filename}"

aws s3 cp --acl public-read "${filepath}" "s3://${bucket}/${filename}"

echo "${url}" | pbcopy
echo "Image url copied to clipboard"
