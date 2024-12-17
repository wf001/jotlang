package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

type TestResult struct {
	Input    string
	Expected string
	Actual   string
	Passed   bool
	Error    error
}

var (
	dir         = fmt.Sprintf("generated/test/%d", time.Now().Unix())
	totalTests  int
	passedCount int
	failedCount int
	mu          sync.Mutex // for safely updating counters
)

func buildCompiler() {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "build", "-o", "./generated/test/modo", "./cmd/modo")
	err = cmd.Run()
	if err != nil {
		panic("Failed to compile the compiler.")
	}
	fmt.Println("\033[0;32mcompiled!\033[0m")
}

func assertExec(
	input string,
	expected string,
	testID int,
	results chan<- TestResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	// 出力ディレクトリをテストごとに一意にする
	outputDir := fmt.Sprintf("%s/out-%d", dir, testID)
	cmd := exec.Command("./generated/test/modo", "run", "-o", outputDir, "--exec", input)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	actualOutput := stdout.String()
	passed := err == nil && actualOutput == fmt.Sprintf("%s\n", expected)

	results <- TestResult{
		Input:    input,
		Expected: expected,
		Actual:   actualOutput,
		Passed:   passed,
		Error:    err,
	}
}

func runTests() {
	fmt.Println("== operation ===")
	testCases := []struct {
		input    string
		expected string
	}{
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 4 13))))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn (+ 1 2 3 4 10))))", "20\n"},
		{"(def main (fn [] (prn (+ 1 2 3 4 5 20))))", "35\n"},
		{"(def main (fn [] (prn (+ 1 2 (+ 3 4)))))", "10\n"},
		{"(def main (fn [] (prn (= 5 (+ 3 2)))))", "1\n"},
	}

	var wg sync.WaitGroup
	results := make(chan TestResult, len(testCases))

	for i, tc := range testCases {
		wg.Add(1)
		go assertExec(tc.input, tc.expected, i, results, &wg)
	}

	wg.Wait()
	close(results)

	for result := range results {
		mu.Lock()
		totalTests++
		if result.Passed {
			passedCount++
			fmt.Printf("%s => %s \033[0;32mOK\033[0m\n", result.Input, result.Actual)
		} else {
			failedCount++
			fmt.Printf("%s => \033[0;31m%s expected, but got %s\033[0m\n", result.Input, result.Expected, result.Actual)
		}
		mu.Unlock()
	}
}

func summary() {
	var status string
	if failedCount > 0 {
		status = "\033[0;31mNG\033[0m"
	} else {
		status = "\033[0;32mOK\033[0m"
	}

	fmt.Println("\n------------------------")
	fmt.Printf(
		"summary: %s, total: %d, passed: %d, failed: %d\n",
		status,
		totalTests,
		passedCount,
		failedCount,
	)
	fmt.Println("------------------------")
}

func main() {
	buildCompiler()
	runTests()
	summary()
}
