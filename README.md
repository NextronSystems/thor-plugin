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
Golang file that defines a package `main`. Additionally, the archive may contain any number of
additional Golang files. If the plugin uses external libraries apart from the standard library,
these must be included in a `vendor` directory in the ZIP archive (see `go mod vendor`).

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

### Examples

See the `examples/` folder for several examples and ideas on how to write THOR plugins.
