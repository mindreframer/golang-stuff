#!/bin/bash

# Run ab with different parameters.
#
# The `nf` app can be found at https://github.com/cloudfoundry/stac2.

NF_URL=http://nf.my.cf.deployment.com/random-data
AB_HOST=root@host.close.to.the.router

C="50 500 2000 5000"
N="50000"
K="1k 16k"

function permute() {
  for c in $C
  do
    for n in $N
    do
      for k in $K
      do
        $@ $c $n $k
      done
    done
  done
}

function run() {
  name=$1

  shift

  c=$1
  n=$2
  k=$3

  AB_COMMAND="ulimit -n $((10 * $c)); ab -v 4 -n $n -c $c $NF_URL?k=$k"

  echo $AB_COMMAND

  ssh $AB_HOST $AB_COMMAND > $name-c$c-n$n-k$k

  echo --- Sleeping some time
  sleep 10
}

path=$(date "+%Y%m%d-%H:%M:%S")

mkdir -p $path

pushd $path

permute run $1
