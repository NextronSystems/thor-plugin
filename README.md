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

Each plugin must define an `Init(thor.Configuration, thor.Logger, thor.RegisterActions)` function.
This function is called on THOR startup and allows plugins to define the conditions when a plugin
should be notified, i.e., register _hooks_.

Hooks are invoked during the scan whenever something is scanned that fulfills the conditions
specified for the hook. In the context of the hook, plugins have access to the data of the scanned
element. There, plugins can perform further analysis on the data, or interact with the running THOR
scan by, e.g., logging a finding, logging an informational message, or manipulate further THOR
actions on the scanned element. Refer to the available hooks in `RegisterActions` in the `thor`
package for more information.

> Note: Plugins are interpreted by [yaegi](https://github.com/traefik/yaegi). While _yaegi_ tries
> to support the Go specification completely for the latest two major Go versions, there are some
> limitations. For instance, plugins cannot use `unsafe` and `syscall` packages from the standard
> library. Refer to _ yaegi_'s documentation for more information.

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

### Examples

See the `examples/` folder for several examples and ideas on how to write THOR plugins.
