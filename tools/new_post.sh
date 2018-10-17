#!/usr/bin/env bash

id=$(uuidgen | tr "[:upper:]" "[:lower:]")
name=${id}.md
prefix=post

hugo new ${prefix}/${name}

vi content/${prefix}/${name}
