package main

import (
	"fmt"

	"github.com/Kqzz/MCsniperGO/pkg/log"
	"github.com/Kqzz/MCsniperGO/pkg/mc"
	"github.com/Kqzz/MCsniperGO/pkg/parser"
)

func getAccounts(giftCodePath string, microsoftPath string) ([]*mc.MCaccount, error) {
	giftCodeLines, _ := parser.ReadLines(giftCodePath)
	microsoftLines, _ := parser.ReadLines(microsoftPath)

	gcs, parseErrors := parser.ParseAccounts(giftCodeLines, mc.MsGc)

	for _, er := range parseErrors {
		if er == nil {
			continue
		}
		log.Log("err", "%v", er)
	}
	microsofts, msParseErrors := parser.ParseAccounts(microsoftLines, mc.Ms)

	for _, er := range msParseErrors {
		if er == nil {
			continue
		}
		log.Log("err", "%v", er)
	}

	accounts := append(gcs, microsofts...)

	if len(accounts) == 0 {
		return accounts, fmt.Errorf("no accounts found in: gc.txt, ms.txt, gp.txt")
	}

	return accounts, nil
}
