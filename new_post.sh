#!/usr/bin/env bash

name=$1.md
prefix=post

# 1. Create the post
hugo new ${prefix}/${name}

# 2. Cat the template 
tail -n 3 content/post/aroad.md >> content/${prefix}/${name}
