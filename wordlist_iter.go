package main

// IterProvider impl and utils for wordlist mode

import (
	"bufio"
	"log/slog"
	"os"
	"strings"

	"github.com/pterm/pterm"
)

type WordlistIter struct {
	cachedPwCount uint64
}

func (iter *WordlistIter) GetPasswordCount() uint64 {
	count := iter.cachedPwCount
	if count != 0 {
		return count
	}

	count = countWordsInFile(*WordlistPath)
	iter.cachedPwCount = count
	return count
}

func (iter *WordlistIter) IterPasswords() func(func(string) bool) {
	return func(yield func(string) bool) {
		f, err := os.Open(*WordlistPath)
		if err != nil {
			pterm.Error.Println(err)
			os.Exit(1)
		}
		defer func() { _ = f.Close() }()

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			word := scanner.Text()
			if !lineIsWord(word) {
				continue
			}

			n := len(word)
			if n > 8 {
				n = 8
			}

			if !yield(word[:n]) {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			slog.Error("error reading from wordlist", "err", err)
		}
	}
}

func lineIsWord(line string) bool {
	return strings.TrimSpace(line) != ""
}

func countWordsInFile(filePath string) (out uint64) {
	f, err := os.Open(filePath)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if lineIsWord(scanner.Text()) {
			out += 1
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("error reading from wordlist", "err", err)
	}

	return out
}
