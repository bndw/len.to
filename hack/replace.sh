# Copyright 2022 Ben Woodward. All rights reserved.
# Use of this source code is governed by a GPL style
# license that can be found in the LICENSE file.
#
#!/usr/bin/env bash
#
# Replaces an image in the len.to s3 bucket. Depends on exiftool
set -e

S3BUCKET=ginput

if [[ "$#" -ne 2 ]] ; then
  cat << EOF

Replaces a file in the len.to s3 bucket.

Usage: $0 <path_to_image> <target_img>
e.g.: $0 /my/new/image.jpg existing.jpg
EOF
  exit 1
fi

if ! [[ -x "$(command -v exiftool)" ]] ; then
  echo 'Error: must install exiftool first.' >&2
  exit 1
fi

filepath=$1
target=$2
filename=$(basename "$1")
extension="${filename##*.}"

# Strip all metadata from the image
exiftool -all= $1

# Upload the image
aws s3 cp --acl public-read "${filepath}" "s3://${S3BUCKET}/${target}"
