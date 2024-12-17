package main

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"os"
	"os/exec"
	"strings"
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
	maxWorkers  int64      = 50
)

func escapeControlChars(s string) string {
	// 改行やタブなどの制御文字をエスケープ
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return strings.TrimRight(s, "\\n")
}

func buildCompiler() {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "build", "-o", "./generated/test/modo", "./cmd/modo")
	if err := cmd.Run(); err != nil {
		panic("Failed to compile the compiler.")
	}
	fmt.Println("\033[0;32mcompiled!\033[0m")
}

func assertExec(
	ctx context.Context,
	input, expected string,
	testID int,
	results chan<- TestResult,
	sem *semaphore.Weighted,
) {
	defer sem.Release(1) // リソース解放

	// 出力ディレクトリをテストごとに一意にする
	outputDir := fmt.Sprintf("%s/out-%d", dir, testID)
	cmd := exec.CommandContext(
		ctx,
		"./generated/test/modo",
		"run",
		"-o",
		outputDir,
		"--exec",
		input,
	)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	actualOutput := escapeControlChars(stdout.String())
	expectedOutput := escapeControlChars(expected)

	passed := err == nil && actualOutput == expectedOutput

	results <- TestResult{
		Input:    input,
		Expected: expectedOutput,
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
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},

		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},

		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},

		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},

		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},

		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
		{"(def main (fn [] (prn (+ 1 2 3))))", "6\n"},
		{"(def main (fn [] (prn 17)))", "17\n"},
	}

	results := make(chan TestResult, len(testCases))
	sem := semaphore.NewWeighted(maxWorkers)
	ctx := context.Background()

	for i, tc := range testCases {
		if err := sem.Acquire(ctx, 1); err != nil {
			fmt.Printf("Failed to acquire semaphore: %v\n", err)
			break
		}

		go func(ctx context.Context, input, expected string, testID int) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered from panic in test %d: %v\n", testID, r)
				}
			}()
			assertExec(ctx, input, expected, testID, results, sem)
		}(ctx, tc.input, tc.expected, i)
	}

	// 全ての Goroutine の終了を待つ
	if err := sem.Acquire(ctx, maxWorkers); err != nil {
		fmt.Printf("Failed to acquire semaphore during finalization: %v\n", err)
	}

	close(results)

	for result := range results {
		mu.Lock()
		totalTests++
		if result.Passed {
			passedCount++
			fmt.Printf("%s => %s \033[0;32mOK\033[0m\n", result.Input, result.Actual)
		} else {
			failedCount++
			fmt.Printf("%s => \033[0;31mExpected: %s, but got: %s\033[0m\n", result.Input, result.Expected, result.Actual)
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
