#!/bin/bash

. ./script/test.sh

testfile(){
  echo "== 1.modo ==="
  assertfile "1.modo" "ans 0\\\ndone\\\n"

}

build-compiler
testfile
summary

exit $code
