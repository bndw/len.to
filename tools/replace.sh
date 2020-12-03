#!/usr/bin/env bash
#
# Replaces an image in the len.to s3 bucket. Depends on exiftool
set -e

if [[ "$#" -ne 1 ]] ; then
  cat << EOF

Replaces a file in the len.to s3 bucket.

Usage: $0 <path_to_image>
EOF
  exit 1
fi

if ! [[ -x "$(command -v exiftool)" ]] ; then
  echo 'Error: must install exiftool first.' >&2
  exit 1
fi

bucket=ginput
filepath=$1
filename=$(basename "$1")
extension="${filename##*.}"

new_filename="${filename}.${extension}"

url="https://d17enza3bfujl8.cloudfront.net/${new_filename}"

# Strip all metadata from the image
exiftool -all= $1

# Upload the image
aws s3 cp --acl public-read "${filepath}" "s3://${bucket}/${new_filename}"
