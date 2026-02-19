## THOR Plugin System

Starting with THOR 11, THOR supports Plugins. THOR Plugins give a quick option to extend THOR with
your own, custom features, e.g.:

* Parse a file format that THOR does not (yet) support
* Check more complex conditions that cannot be written as custom IOCs or rules
* Post-processing: Extend THOR output in custom, user-defined ways

In a nutshell, THOR Plugins are ZIP archives containing Golang code that is executed by THOR during
a scan. Plugins register _hooks_ that are called during the scan and perform custom actions in
there.

### Using a Plugin

Plugins are ZIP files placed in a `plugins/` folder in your THOR directory. Each file is
interpreted as a separate plugin.

> Warning: Plugins contain executable code that is run by THOR. For this reason, never run any
> plugins that do not come from a trusted source.

### Writing a Plugin

Plugins are written in Golang and communicate with THOR via an interface defined in
`thorplugin.go`. They are packaged as ZIP archives and placed in the `plugins/` directory.

The ZIP archive's content is similar to an independent Golang package. The archive must contain a
Golang file that defines a package `main`. Additionally, the archive may contain:
* Any number of additional Golang files
* A `metadata.yml` file with information about the plugin (see below)
* A `vendor` directory in case the plugin uses external libraries apart from the standard library (see `go mod vendor`)

> Note: Plugins are interpreted by [yaegi](https://github.com/traefik/yaegi). While _yaegi_ tries
> to support the Go specification completely for the latest two major Go versions, there are some
> limitations. For instance, plugins cannot use `unsafe` and `syscall` packages from the standard
> library. Refer to _yaegi_'s documentation for more information.

Each plugin must define an `Init(thor.Configuration, thor.Logger, thor.RegisterActions)` function.
This function is called on THOR startup and allows plugins to define the conditions when a plugin
should be notified, i.e., register _hooks_.

Hooks are invoked during the scan whenever something is scanned that fulfills the conditions
specified for the hook.

In the context of the hook, plugins have access to the data of the scanned
element and can perform further analysis on the data, or interact with the running THOR
scan. This could be e.g. logging a finding, logging an informational message, or manipulate further THOR
actions on the scanned element.

#### Tag hooks

The most common hook is `AddRuleHook`, which is based on _tags_. All
[signatures](https://thor-manual.nextron-systems.com/en/v11/signatures/index.html) can have tags;
whenever a signature with a specific tag matches on an object, the hook is triggered.

This is used for e.g. parsers; define a YARA rule that matches on the file format you want
to parse, give it a unique tag, and have your parser hook trigger on this tag.

> YARA rules for file formats should usually be
> [meta rules](https://thor-manual.nextron-systems.com/en/v11/signatures/yara.html#specific-yara-rules)
> to ensure that they are applied to all files, i.e., use `TypeMeta` as `YaraRuleType`.

```go
func Init(config thor.Configuration, logger thor.Logger, actions thor.RegisterActions) {
    actions.AddYaraRule(thor.TypeMeta, `
rule DetectMyFormat: MYTAG {
    meta:
        score = 0
    strings:
        $magicheader = "deadbeef"
    condition: $magicheader at 0
}`)
    actions.AddRuleHook("MYTAG", parseMyFormat)
}

func parseMyFormat(scanner thor.Scanner, object thor.MatchingObject) {
    // Analyze the object here...
}
```

#### Post processing hooks

Post processing hooks can be registered with `AddPostProcessingHook` and are invoked for every element.
They have access to the complete report that THOR generates for this object, including all signatures that matched.

They are usually used for logging; e.g. you could send all findings to your SIEM in a custom format.

#### Metadata

Plugins may contain a `metadata.yml` file in the root of the ZIP archive. This file contains
metadata about the plugin, such as the plugin's name, version, and a description. \
THOR reads this file and displays the information in the THOR log when the plugin is loaded.

The `metadata.yml` file must be a valid YAML file and may contain the following fields:

* `name`: The name of the plugin. This field is mandatory and will be used in the THOR log for output from the plugin.
* `version`: The version of the plugin. This field is optional and may be used to track the plugin's version.
* `description`: A description of the plugin. This field is optional and may be used to describe the plugin's purpose.
* `author`: The author of the plugin. This field is optional and may be used to credit the plugin's author.
* `requires_thor`: The minimum THOR version required to run the plugin.
  This field is optional.
  It must be a valid semantic version string, prefixed with `v` (e.g., `v11.0.0`). \
  If the THOR version is lower than the specified version, THOR will not load the plugin and log an error instead.
* `link`: A URL to the plugin's source code or documentation. This field is optional.
* `build_tags`: A list of build tags that are applied when loading the plugin. This field is optional.
  It may be used to specify build tags that are required for the plugin to work correctly.

### Examples

See [THOR Plugins](https://github.com/thor-plugins/) for existing plugins, including several examples.
There is also a [plugin template](https://github.com/thor-plugins/template) available that can be used as a starting point for new plugins.
