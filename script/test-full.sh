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
  actual_output=$(./generated/test/modo run -o "$dir/out" --exec "$input" |gsed ':a;N;$!ba;s/\n/\\\\n/g')
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
  # multiline
  echo "== multi line ==="
  assertexec '(def main (fn [] (prn (+ 1 2)) (prn (+ 3 4))))' "3\\\n7\\\n"
  assertexec '(def x 4) (def main (fn [] (let [y 2 z (+ x 3)] (prn (+ x z)) (prn (+ x y)))))' "11\\\n6\\\n"

  # operator
  echo "== operation ==="
  assertexec '(def main (fn [] (prn 17)))' "17\\\n"
  assertexec '(def main (fn [] (prn (+ 4 13))))' "17\\\n"
  assertexec '(def main (fn [] (prn (+ 1 2 3))))' "6\\\n"
  assertexec '(def main (fn [] (prn (+ 1 2 3 4 10))))' "20\\\n"
  assertexec '(def main (fn [] (prn (+ 1 2 3 4 5 20))))' "35\\\n"
  assertexec '(def main (fn [] (prn (+ 1 2 (+ 3 4)))))' "10\\\n"
  assertexec '(def main (fn [] (prn (+ (+ 1 2) (+ 3 4)))))' "10\\\n"
  assertexec '(def main (fn [] (prn (+ (+ 1 2) (+ (+ 9 5) 4)))))' "21\\\n"
  assertexec '(def main (fn [] (prn (+ 1 (+ 3 2) (+ (+ 9 4 5) 7 8)))))' "39\\\n"
  assertexec '(def main (fn [] (prn (= 5 (+ 3 2)))))' "1\\\n"
  assertexec '(def main (fn [] (prn (= (+ 4 3) (+ 3 2)))))' "0\\\n"
  # as often failed
  assertexec '(def main (fn [] (prn (= (+ 4 -3) (+ 3 -2)))))' "1\\\n"

  # variable
  echo "== variable ==="
  assertexec '(def x 1) (def main (fn [] (prn (+ x 2))))' "3\\\n"
  assertexec '(def x 1) (def main (fn [] (prn (+ 2 x))))' "3\\\n"
  assertexec '(def x 1) (def main (fn [] (prn (+ x (+ 2 3)))))' "6\\\n"
  assertexec '(def x 1) (def main (fn [] (prn (+ x (+ 2 3) 4))))' "10\\\n"
  assertexec '(def x 1) (def y 2) (def main (fn [] (prn (+ x y))))' "3\\\n"
  assertexec '(def x 1) (def main (fn [] (prn (+ x x))))' "2\\\n"

  assertexec '(def x 1) (def main (fn [] (prn (= x 2))))' "0\\\n"
  assertexec '(def x 1) (def main (fn [] (prn (= x 1))))' "1\\\n"
  assertexec '(def x 1) (def y 2) (def main (fn [] (prn (= x y))))' "0\\\n"

  # binding
  echo "== binding ==="
  assertexec '(def main (fn [] (let [x 1] (prn x))))' "1\\\n"
  assertexec '(def x 4) (def main (fn [] (let [y 2] (prn (+ x y)))))' "6\\\n"
  assertexec '(def x 4) (def main (fn [] (let [y 2 z (+ y 3)] (prn (+ x z)))))' "9\\\n"

  # if
  echo "== if ==="
  assertexec '(def main (fn [] (if (= 1 1) (prn 11) (prn 12))))' "11\\\n"
  assertexec '(def main (fn [] (if (= 1 2) (prn 11) (prn 12))))' "12\\\n"

  assertexec '(def main (fn [] (if (= 1 1) (prn 11) (if (= 1 1) (prn 12) (prn 13)))))' "11\\\n"
  assertexec '(def main (fn [] (if (= 1 1) (prn 11) (if (= 1 2) (prn 12) (prn 13)))))' "11\\\n"
  assertexec '(def main (fn [] (if (= 1 2) (prn 11) (if (= 1 1) (prn 12) (prn 13)))))' "12\\\n"
  assertexec '(def main (fn [] (if (= 1 2) (prn 11) (if (= 1 2) (prn 12) (prn 13)))))' "13\\\n"

  assertexec '(def main (fn [] (if (= 1 1) (if (= 1 1) (prn 11) (prn 12)) (prn 13))))' "11\\\n"
  assertexec '(def main (fn [] (if (= 1 1) (if (= 1 2) (prn 11) (prn 12)) (prn 13))))' "12\\\n"
  assertexec '(def main (fn [] (if (= 1 2) (if (= 1 1) (prn 11) (prn 12)) (prn 13))))' "13\\\n"
  assertexec '(def main (fn [] (if (= 1 2) (if (= 1 2) (prn 11) (prn 12)) (prn 13))))' "13\\\n"

  assertexec '(def main (fn [] (if (= 1 2) (prn 11) (if (= 1 1) (if (= 1 2) (prn 12) (prn 13)) (prn 14) ))))' "13\\\n"

  # string
  echo "== string ==="
  assertexec '(def main (fn [] (prn "hello")))' "hello\\\n"
  assertexec '(def main (fn [] (let [s "hello"] (prn s))))' "hello\\\n"

  # prn
  assertexec '(def main (fn [] (prn 1 2 "hello")))' "1 2 hello\\\n"

}

build-compiler
testexec
summary

exit $code

