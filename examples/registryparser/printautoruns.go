package main

import (
	"path/filepath"

	"github.com/NextronSystems/jsonlog/thorlog/v3"
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
	actions.AddRuleHook("RUNKEY", func(scanner thor.Scanner, object thor.MatchingObject) {
		registryValue, isRegistryValue := object.Object.(*thorlog.RegistryValue)
		if !isRegistryValue {
			return
		}
		valueName := filepath.Base(registryValue.Key)
		logger.Info("Found autorun entry in registry", "value", valueName, "command", registryValue.ParsedValue)
	})
	logger.Info("PrintAutoruns plugin loaded!")
}
