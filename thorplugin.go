package thor

import (
	"io"
)

// Initial entry point for THOR plugins is:
// package main
// func Init(Configuration, Logger, RegisterActions)
// All plugins must provide this function.

type PluginInitializeFunction func(Configuration, Logger, RegisterActions)

type Configuration struct {
	MaxFileSize uint64
}

type RegisterActions interface {
	AddYaraRule(ruletype RuleType, rule string)
	AddYaraRuleHook(tag string, callback RuleMatchedCallback)
}

type RuleType int

const (
	TypeMeta RuleType = iota
	TypeKeyword
	TypeDefault
	TypeRegistry
	TypeLog
	TypeProcess
)

type RuleMatchedCallback func(scanner Scanner, object MatchingObject)

type MatchingObject struct {
	Reader ObjectReader
	Path   string
}

type ObjectReader interface {
	io.ReaderAt
	io.ReadSeeker
	Size() int64
}

type Scanner interface {
	ScanString(data string)
	ScanFile(name string, data []byte)

	Logger
}

type Logger interface {
	Log(score int64, text string)
	Info(text string)
	Debug(text string)
	Error(text string)
}
