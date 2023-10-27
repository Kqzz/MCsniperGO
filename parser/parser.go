package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Kqzz/MCsniperGO/mc"
)

func ParseAccounts(accs []string, accType mc.AccType) ([]*mc.MCaccount, []error) {
	parsed, errs := []*mc.MCaccount{}, []error{}
	for i, l := range accs {

		if len(l) > 0 && l[0] == '#' { // commented
			continue
		}

		if len(l) > 200 &&
			!strings.Contains(l, ":") &&
			strings.HasPrefix(l, "eyJ") { // bearer token
			acc := &mc.MCaccount{Email: l[40:50], Type: accType, Bearer: l}
			acc.DefaultFastHttpHandler()
			parsed = append(parsed, acc)
			continue
		}

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

		acc.DefaultFastHttpHandler()

		parsed = append(parsed, acc)

	}
	return parsed, errs
}

func ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)

	if err != nil {
		os.Create(filename)
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
