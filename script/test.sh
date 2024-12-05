#!/bin/bash

dir="generated/test/$(date +%s)"
testdatadir="./testdata/modo"
msg="\033[0;32mOK\033[0m"
code=0
total_tests=0
passed_count=0
failed_count=0

assertexec() {
  input="$1"
  expected="$2"
  assert "$input" "$expected" 
}

assertfile() {
  input="$1"
  expected=$(cat "$testdatadir/$2")
  assert "$input" "$expected" 
}

assert() {
  input="$1"
  expected="$2"
  actual_output=$(./generated/test/modo run -o "$dir/out" --exec "$input")
  actual_exit_code="$?"

  ((total_tests++))

  diff_result=$(diff <(echo "$expected") <(echo "$actual_output"))
  if [ "$actual_exit_code" -eq 0 ] && [ -z "$diff_result" ]; then
    ((passed_count++))
    echo -e "$input => $actual_output \033[0;32mOK\033[0m"
  else
    ((failed_count++))
    echo -e "$input => \033[0;31m$expected expected, but got $actual_output\033[0m"
    msg="\033[0;31mNG\033[0m"
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

testexec(){
  # operator
  assertexec '(def main (fn [] (prn 17)))' 17
  assertexec '(def main (fn [] (prn (+ 4 13))))' 17
  assertexec '(def main (fn [] (prn (+ 1 2 3))))' 6
  assertexec '(def main (fn [] (prn (+ 1 2 3 4 10))))' 20
  assertexec '(def main (fn [] (prn (+ 1 2 3 4 5 20))))' 35
  assertexec '(def main (fn [] (prn (+ 1 2 (+ 3 4)))))' 10
  assertexec '(def main (fn [] (prn (+ (+ 1 2) (+ 3 4)))))' 10
  assertexec '(def main (fn [] (prn (+ (+ 1 2) (+ (+ 9 5) 4)))))' 21
  assertexec '(def main (fn [] (prn (+ 1 (+ 3 2) (+ (+ 9 4 5) 7 8)))))' 39
  assertexec '(def main (fn [] (prn (= 5 (+ 3 2)))))' 1
  assertexec '(def main (fn [] (prn (= (+ 4 3) (+ 3 2)))))' 0

  # variable
  assertexec '(def x 1) (def main (fn [] (prn (+ x 2))))' 3
}
testfile(){
  assertfile '(def main (fn [] (prn (+ 1 2)) (prn (+ 3 4))))' 'SimpleSequentialOutput1' 
}

build-compiler
testexec
# not work
testfile
summary

exit $code

