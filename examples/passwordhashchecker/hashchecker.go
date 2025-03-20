package main

import (
	"bufio"
	"strings"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/thorlog/v3"
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
		var offset int
		for lineReader.Scan() {
			fullLine := lineReader.Text()

			line := strings.Split(fullLine, ":")
			if len(line) != 9 {
				scanner.Error("Corrupt shadow line")
				offset += len(fullLine) + 1
				continue
			}
			user, hash := line[0], line[1]
			if strings.HasPrefix(hash, "$1$") {
				userOffset := uint64(offset)
				hashOffset := uint64(offset + len(user) + 1)
				scanner.AddReason(thorlog.NewReason("User has MD5 hashed password in /etc/shadow file", thorlog.Signature{
					Type:  thorlog.Custom,
					Class: thorlog.ClassInternalHeuristic,
					Score: 60,
				}, thorlog.MatchStrings{
					{
						Match:  thorlog.MatchData{Data: []byte(user)},
						Offset: &userOffset,
						Field:  jsonlog.NewReference(object.Object, &object.Object.(*thorlog.File).Content),
					},
					{
						Match:  thorlog.MatchData{Data: []byte(hash)},
						Offset: &hashOffset,
						Field:  jsonlog.NewReference(object.Object, &object.Object.(*thorlog.File).Content),
					},
				}))
			}
			offset += len(fullLine) + 1
		}
	})
	logger.Info("ShadowParser plugin loaded!")
}
