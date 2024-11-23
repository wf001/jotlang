#!/bin/bash

dir="generated/test/$(date +%s)"
msg="\033[0;32mOK\033[0m"
code=0
total_tests=0
passed_count=0
failed_count=0

assert() {
  expected="$1"
  input="$2"

  # 実行結果を変数に格納
  actual_output=$(./generated/test/modo run -o "$dir/out" --exec "$input")
  actual_exit_code="$?"

  ((total_tests++))

  # 実行結果をdiffで比較
  diff_result=$(diff <(echo "$expected") <(echo "$actual_output"))
  if [ "$actual_exit_code" -eq 0 ] && [ -z "$diff_result" ]; then
    ((passed_count++))
    echo -e "\033[0;32m$input => OK\033[0m"
  else
    ((failed_count++))
    echo -e "\033[0;31m$input => $expected expected, but got $actual_output\033[0m"
  fi
}

build-compiler(){
  mkdir -p "$dir"
  go build -o ./generated/test/modo ./cmd/modo 
  echo -e "\033[0;32mcompiled!\033[0m"
}

summary(){
  echo -e "\n------------------------"
  echo -e "summary: $msg, total: $total_tests, passed: $passed_count, failed: $failed_count"
  echo -e "------------------------"
}

testit(){
  assert 17 '(prn 17)'
  assert 17 '(prn (+ 4 13))'
  assert 6 '(prn (+ 1 2 3))'
  assert 20 '(prn (+ 1 2 3 4 10))'
  assert 35 '(prn (+ 1 2 3 4 5 20))'
  assert 10 '(prn (+ 1 2 (+ 3 4)))'
  assert 10 '(prn (+ (+ 1 2) (+ 3 4)))'
  assert 21 '(prn (+ (+ 1 2) (+ (+ 9 5) 4)))'
  assert 39 '(prn (+ 1 (+ 3 2) (+ (+ 9 4 5) 7 8)))'
  assert 1 '(prn (= 5 (+ 3 2)))'
  assert 0 '(prn (= (+ 4 3) (+ 3 2)))'
}

build-compiler
testit
summary

exit $code

