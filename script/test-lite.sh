#!/bin/bash

. ./script/test.sh

testfile(){
  echo "== 1.modo ==="
  assertfile "1.modo" "ans 456\\\ndone\\\n"
  echo "== 2.modo ==="
  assertfile "2.modo" "Another\\\n"
  echo "== 3.modo ==="
  assertfile "3.modo" "45\\\n"

}

build-compiler
testfile
summary

exit $code
