package main

import (
	"fmt"

	"github.com/Kqzz/MCsniperGO/log"
	"github.com/Kqzz/MCsniperGO/pkg/mc"
	"github.com/Kqzz/MCsniperGO/pkg/parser"
)

func getAccounts(giftCodePath string, gamepassPath string, microsoftPath string) ([]*mc.MCaccount, error) {
	giftCodeLines, _ := parser.ReadLines(giftCodePath)
	gamepassLines, _ := parser.ReadLines(gamepassPath)
	microsoftLines, _ := parser.ReadLines(microsoftPath)

	gcs, parseErrors := parser.ParseAccounts(giftCodeLines, mc.MsPr)

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

	gamepasses, gpParseErrors := parser.ParseAccounts(gamepassLines, mc.MsGp)

	for _, er := range gpParseErrors {
		if er == nil {
			continue
		}

	}

	accounts := append(gcs, microsofts...)
	accounts = append(accounts, gamepasses...)

	if len(accounts) == 0 {
		return accounts, fmt.Errorf("no accounts found in: gc.txt, ms.txt, gp.txt")
	}

	return accounts, nil
}
