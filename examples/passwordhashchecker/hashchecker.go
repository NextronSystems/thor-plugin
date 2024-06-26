package main

import (
	"bufio"
	"strings"

	"github.com/NextronSystems/thor-plugin"
)

func Init(config thor.Configuration, logger thor.Logger, actions thor.RegisterActions) {
	actions.AddYaraRule(thor.TypeMeta, `
rule Shadow: SHADOWFILE {
    meta:
        score = 0
    condition: filepath == "/etc" and filename == "shadow"
}`)
	actions.AddRuleHook("SHADOWFILE", func(scanner thor.Scanner, object thor.MatchingObject) {
		lineReader := bufio.NewScanner(object.Content)
		for lineReader.Scan() {
			fullLine := lineReader.Text()

			line := strings.Split(fullLine, ":")
			if len(line) != 9 {
				scanner.Error("Corrupt shadow line")
				continue
			}
			user, hash := line[0], line[1]
			if strings.HasPrefix(hash, "$1$") {
				scanner.Log("Unsafe user details found", "User has MD5 hashed password in /etc/shadow file", 60, "user", user)
			}
		}
	})
	logger.Info("ShadowParser plugin loaded!")
}
