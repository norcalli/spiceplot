#!/usr/bin/env bash
name=$1
shift 1
ngspice -b $name.net -r $name.raw >/dev/null && {
  ./spiceplot "$@" $name.raw $name.svg
  rm $name.raw
}