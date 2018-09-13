#!/bin/bash

for fn in content/post/* ; do
  echo $fn
  #sed 's/.\{29\}$/"/' $fn
  sed 's/^title*.{12}.\{19}/"/' $fn
done
