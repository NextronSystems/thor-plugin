package main

import (
	"bufio"
	"strings"

	"github.com/NextronSystems/thor-plugin"
)

func Init(config thor.Configuration, logger thor.Logger, actions thor.RegisterActions) {
	actions.AddYaraRule(thor.TypeRegistry, `
rule RunKey: RUNKEY {
    meta:
        score = 0
	strings:
		$s1 = "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run" nocase
    condition:
		1 of them
}`)
	actions.AddYaraRuleHook("RUNKEY", func(scanner thor.Scanner, object thor.MatchingObject) {
		// See https://thor-manual.nextron-systems.com/en/latest/usage/custom-signatures.html#thor-yara-rules-for-registry-detection
		// for the format of the data we receive from the registry
		lineReader := bufio.NewScanner(object.Reader)
		for lineReader.Scan() {
			splitLine := strings.SplitN(lineReader.Text(), ";", 3)
			logger.Info("Found autorun entry in registry", "value", splitLine[1], "command", splitLine[2])
		}
	})
	logger.Info("PrintAutoruns plugin loaded!")
}
