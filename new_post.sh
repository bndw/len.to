#!/usr/bin/env bash

id=$(uuidgen | tr "[:upper:]" "[:lower:]")
name=${id}.md
prefix=post

# 1. Create the post
hugo new ${prefix}/${name}

# 2. Cat the template 
tail -n 4 content/post/olives_sun.md >> content/${prefix}/${name}
vi content/${prefix}/${name}
