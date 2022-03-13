package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Kqzz/MCsniperGO/mc"
)

func parseAccounts(accs []string) ([]*mc.MCaccount, []error) {
	parsed, errs := []*mc.MCaccount{}, []error{}
	for i, l := range accs {
		s := strings.Split(l, ":")

		if len(s) == 0 {
			continue
		}

		acc := &mc.MCaccount{}

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

	rd := bufio.NewReader(file)

	lines := []string{}

	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			return []string{}, err
		}

		lines = append(lines, line)
	}

	return lines, nil
}
