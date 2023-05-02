package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Kqzz/MCsniperGO/mc"
)

func parseAccounts(accs []string, accType mc.AccType) ([]*mc.MCaccount, []error) {
	parsed, errs := []*mc.MCaccount{}, []error{}
	for i, l := range accs {
		s := strings.Split(l, ":")

		if len(s) == 0 {
			continue
		}

		acc := &mc.MCaccount{Type: accType}

		if len(s) >= 2 {
			acc.Email = s[0]
			acc.Password = s[1]
		} else {
			errs = append(errs, fmt.Errorf("invalid split count on line %v", i))
			continue
		}

		parsed = append(parsed, acc)

	}
	return parsed, errs
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)

	if err != nil {
		return []string{}, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}
