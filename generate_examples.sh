#!/bin/bash

if [ $# -le 1 ]; then
  echo "No arguments supplied"
fi

EXAMPLE="name: A GET request {NR}
url: 'https://httpbin.org/anything'
method: GET
headers:
  X-Foo-Bar: SomeValue
body: >
  {
    "id": 1,
    "name": "Jane"
  }
"

for i in $(seq 1 $1); do
  nr="$i"
  if (($i < 10)); then
    nr="0$i"
  fi
  REQ=$(echo "$EXAMPLE" | sed "s/{NR}/$nr/g")
  echo "$REQ" >$2/"A GET request $nr.yaml"
done
