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
  # arithmetic operator
  echo "== arithmetic operator ==="
  assertexec '(def main ::int (fn [] (prn 17)))' "17\\\n"
  assertexec '(def main ::int (fn [] (prn (+ 4 13))))' "17\\\n"
  assertexec '(def main ::int (fn [] (prn (+ 1 2 3))))' "6\\\n"
  assertexec '(def main ::int (fn [] (prn (+ 1 2 3 4 10))))' "20\\\n"
  assertexec '(def main ::int (fn [] (prn (+ 1 2 3 4 5 20))))' "35\\\n"
  assertexec '(def main ::int (fn [] (prn (+ 1 2 (+ 3 4)))))' "10\\\n"
  assertexec '(def main ::int (fn [] (prn (+ (+ 1 2) (+ 3 4)))))' "10\\\n"
  assertexec '(def main ::int (fn [] (prn (+ (+ 1 2) (+ (+ 9 5) 4)))))' "21\\\n"
  assertexec '(def main ::int (fn [] (prn (+ 1 (+ 3 2) (+ (+ 9 4 5) 7 8)))))' "39\\\n"

  # equality operator
  echo "== equality operator ==="
  assertexec '(def main ::int (fn [] (prn (= 123 123))))' "1\\\n"
  assertexec '(def main ::int (fn [] (prn (= 5 (+ 3 2)))))' "1\\\n"
  assertexec '(def main ::int (fn [] (prn (= (+ 4 3) (+ 3 2)))))' "0\\\n"
  assertexec '(def main ::int (fn [] (prn (= (+ 4 -3) (+ 3 -2)))))' "1\\\n"
  assertexec '(def main ::int (fn [] (prn (= "foo" "bar"))))' "0\\\n"
  assertexec '(def main ::int (fn [] (prn (= "foo" "foo"))))' "1\\\n"

  ## global variable
  echo "== global variable ==="
  assertexec '(def x ::int 3) (def main ::int (fn [] (prn x)))' "3\\\n"
  assertexec '(def x ::string "hello") (def main ::int (fn [] (prn x)))' "hello\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (+ x 2))))' "3\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (+ 2 x))))' "3\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (+ x (+ 2 3)))))' "6\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (+ x (+ 2 3) 4))))' "10\\\n"
  assertexec '(def x ::int 1) (def y ::int 2) (def main ::int (fn [] (prn (+ x y))))' "3\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (+ x x))))' "2\\\n"

  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (= x 2))))' "0\\\n"
  assertexec '(def x ::int 1) (def main ::int (fn [] (prn (= x 1))))' "1\\\n"
  assertexec '(def x ::int 1) (def y ::int 2) (def main ::int (fn [] (prn (= x y))))' "0\\\n"

  ## binded variable
  echo "== binded variable ==="
  assertexec '(def main ::int (fn [] (let [x ::int 1] (prn x))))' "1\\\n"
  assertexec '(def main ::int (fn [] (let [x ::int 1 y ::int (+ x 2)] (prn y))))' "3\\\n"
  assertexec '(def x ::int 4) (def main ::int (fn [] (let [y ::int 2] (prn (+ x y)))))' "6\\\n"
  assertexec '(def x ::int 4) (def main ::int (fn [] (let [y ::int 2 z ::int (+ y 3)] (prn (+ x z)))))' "9\\\n"

  ## if
  echo "== if ==="
  assertexec '(def main ::int (fn [] (if (= 1 1) (prn 11) (prn 12))))' "11\\\n"
  assertexec '(def main ::int (fn [] (if (= 1 2) (prn 11) (prn 12))))' "12\\\n"

  assertexec '(def main ::int (fn [] (if (= 1 1) (prn 11) (if (= 1 1) (prn 12) (prn 13)))))' "11\\\n"
  assertexec '(def main ::int (fn [] (if (= 1 1) (prn 11) (if (= 1 2) (prn 12) (prn 13)))))' "11\\\n"
  assertexec '(def main ::int (fn [] (if (= 1 2) (prn 11) (if (= 1 1) (prn 12) (prn 13)))))' "12\\\n"
  assertexec '(def main ::int (fn [] (if (= 1 2) (prn 11) (if (= 1 2) (prn 12) (prn 13)))))' "13\\\n"

  assertexec '(def main ::int (fn [] (if (= 1 1) (if (= 1 1) (prn 11) (prn 12)) (prn 13))))' "11\\\n"
  assertexec '(def main ::int (fn [] (if (= 1 1) (if (= 1 2) (prn 11) (prn 12)) (prn 13))))' "12\\\n"
  assertexec '(def main ::int (fn [] (if (= 1 2) (if (= 1 1) (prn 11) (prn 12)) (prn 13))))' "13\\\n"
  assertexec '(def main ::int (fn [] (if (= 1 2) (if (= 1 2) (prn 11) (prn 12)) (prn 13))))' "13\\\n"

  assertexec '(def main ::int (fn [] (if (= 1 2) (prn 11) (if (= 1 1) (if (= 1 2) (prn 12) (prn 13)) (prn 14) ))))' "13\\\n"

  ## string type
  echo "== string type ==="
  assertexec '(def main ::int (fn [] (prn "hello")))' "hello\\\n"
  assertexec '(def main ::int (fn [] (let [s ::string "hello"] (prn s))))' "hello\\\n"
  assertexec '(def main ::int (fn [] (let [s ::string "hello" t ::string "world"] (prn s))))' "hello\\\n"
  assertexec '(def main ::int (fn [] (let [s ::string "hello" t ::string "world"] (prn t))))' "world\\\n"

  ## prn multi-line
  echo "== prn multi-line ==="
  assertexec '(def main :: int (fn [] (prn (+ 1 2)) (prn (+ 3 4))))' "3\\\n7\\\n"
  assertexec '(def x :: int 4) (def main ::int (fn [] (let [y ::int 2 z ::int (+ x 3)] (prn (+ x z)) (prn (+ x y)))))' "11\\\n6\\\n"

  ## prn multi-value
  echo "== prn multi-value ==="
  assertexec '(def main :: int (fn [] (prn 1 2 "hello")))' "1 2 hello\\\n"
  assertexec '(def main :: int (fn [] (prn (+ 1 2) (= 2 4) "world")))' "3 0 world\\\n"

  # function calling
  echo "== function calling ==="
  assertexec '(def f :: int => nil (fn [a] (prn a))) (def main ::int (fn [] (f 4)))' "4\\\n"
  assertexec '(def f :: int => string => nil (fn [a s] (prn a) (prn s))) (def main ::int (fn [] (f 123 "hello")))' "123\\\nhello\\\n"
  assertexec '(def f :: int => int => int (fn [a b] (+ a b))) (def main ::int (fn [] (prn (f 1 2))))' "3\\\n"
  assertexec '(def f :: string => string (fn [a] a)) (def main ::int (fn [] (prn (f "hello"))))' "hello\\\n"
  assertexec '(def f :: string => string (fn [a] "modo")) (def main ::int (fn [] (prn (f "hello"))))' "modo\\\n"
  assertexec '(def f :: string => string (fn [a] (let [s ::string "modo"] s))) (def main ::int (fn [] (prn (f "hello"))))' "modo\\\n"

}

build-compiler
testexec
summary

exit $code

