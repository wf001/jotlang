#!/bin/bash

dir="generated/test/$(date +%s)"
msg="OK"
code=0
total_tests=0
passed_count=0
failed_count=0

assert() {
  expected="$1"
  input="$2"

  ./"$dir"/out

  actual="$?"

  ((total_tests++))

  if [ "$actual" = "$expected" ]; then
    ((passed_count++))
    echo -e "\033[0;32m$input => $actual\033[0m"
  else
    ((failed_count++))
    echo -e "\033[0;31m$input => $expected expected, but got $actual\033[0m"
    msg="NG"
    code=-1
  fi
}

build(){
  mkdir -p "$dir"
  go run ./cmd/modo > "$dir"/out.ll
  llc "$dir"/out.ll -o "$dir"/out.s
  clang "$dir"/out.s -o "$dir"/out
}

summary(){
  echo -e "\n------------------------"
  echo -e "summary: $msg, total: $total_tests, passed: $passed_count, failed: $failed_count"
  echo -e "------------------------"
}

testit(){
  assert 5 '5'
}

build
testit
summary

exit $code

