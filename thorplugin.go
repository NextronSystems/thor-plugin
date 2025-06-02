package thor

import (
	"io"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/thorlog/v3"
)

// Initial entry point for THOR plugins is:
// package main
// func Init(Configuration, Logger, RegisterActions)
// All plugins must provide this function.

type PluginInitializeFunction func(Configuration, Logger, RegisterActions)

// Configuration contains information about the parameters THOR was started with.
type Configuration struct {
	MaxFileSize uint64
}

// RegisterActions provides ways to register with THOR during the plugin initialization.
// The provided functions can only be called from the plugin initialization; once the
// plugin initialization is complete, they will no longer have an effect.
type RegisterActions interface {

	// AddYaraRule adds one or multiple YARA rules to THOR's ruleset.
	//
	// This is typically used to register a rule with a special tag that is then used with
	// AddRuleHook.
	AddYaraRule(ruletype YaraRuleType, rule string)

	// AddRuleHook registers a callback for a specific rule tag.
	//
	// Whenever a YARA or sigma rule with this tag matches on any data, the callback
	// is invoked with the data that the rule matched on.
	// The matched data can be a file, registry entry, log entry, or any
	// other kind of data that is scanned by THOR.
	AddRuleHook(tag string, callback RuleMatchedCallback)

	AddPostProcessingHook(callback PostProcessingCallback)
}

// YaraRuleType defines a type of YARA rules within THOR.
// Each rule type is applied to different type of data.
type YaraRuleType int

const (
	// TypeMeta rules are applied to all files, however, they can only access
	// the first 2048 bytes of each file and the THOR external variables.
	TypeMeta YaraRuleType = iota

	// TypeKeyword are applied to all elements except for files.
	TypeKeyword

	// TypeDefault are applied to files where
	// a deep scan was started (typically decided by magic header, extension and file size).
	TypeDefault

	// TypeRegistry rules are applied exclusively to registry data.
	TypeRegistry

	// TypeLog rules are applied exclusively to log data (log files or event logs).
	TypeLog

	// TypeProcess rules are applied to scanned processes.
	TypeProcess
)

// RuleMatchedCallback describes a callback for matched rules.
type RuleMatchedCallback func(scanner Scanner, object MatchingObject)

// PostProcessingCallback describes a callback for actions on a fully scanned object.
type PostProcessingCallback func(logger Logger, object MatchedObject)

// MatchingObject describes an object that a rule matched on.
type MatchingObject struct {
	// Object is the full description of the object that the rule matched on.
	Object jsonlog.Object
	// Reader provides access to the content of the object that the rule matched on. The content will be empty for all objects except for files and processes.
	Content ObjectReader
}

type MatchedObject struct {
	Finding *thorlog.Finding
	Content ObjectReader
}

type ObjectReader interface {
	io.ReaderAt
	io.ReadSeeker
	Size() int64
}

type KeyValuePair struct {
	Key   string
	Value string
}

// Scanner provides methods within a RuleMatchedCallback to scan further data (typically data extracted
// from the MatchingObject passed to the callback).
//
// Each scanner instance is only valid for the duration of the callback.
type Scanner interface {
	// ScanString scans a passed string with filename IOCs, keyword YARA rules, and Sigma rules.
	ScanString(data string)

	// ScanFile scans a passed file as if it was found on the file system.
	// unpackMethod should contain the method by which this file was extracted from the
	// MatchingObject. It is used for the file's unpack_source and unpack_parent YARA external variables
	// and should by convention be an upper case word, e.g. ZIP or RAR.
	ScanFile(name string, data []byte, unpackMethod string)

	// ScanStructuredData scans a set of key/value pairs with filename IOCs, keyword YARA rules,
	// and Sigma rules.
	ScanStructuredData(data []KeyValuePair)

	// AddReason adds a reason to the element passed to the callback.
	AddReason(reason thorlog.Reason)

	Logger
}

type Logger interface {
	// Info logs an informational message.
	//
	// kv needs to be a number of key / value pairs, where each key must be a
	// string, ordered as key1, value1, key2, value2, ...
	Info(text string, kv ...any)

	// Debug logs a debug level message.
	//
	// kv needs to be a number of key / value pairs, where each key must be a
	// string, ordered as key1, value1, key2, value2, ...
	Debug(text string, kv ...any)

	// Error logs an error message.
	//
	// kv needs to be a number of key / value pairs, where each key must be a
	// string, ordered as key1, value1, key2, value2, ...
	Error(text string, kv ...any)
}
