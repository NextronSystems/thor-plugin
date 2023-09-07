## THOR Plugin System

Starting with THOR 10.8, THOR supports Plugins. THOR Plugins give a quick option to extend THOR with your own,
custom features by executing your own scanning or parsing logic within THOR.

### Introduction

Plugins are Golang files placed in a `plugins/` folder in your THOR directory. Each file is interpreted as a separate plugin.

Each plugin must define an `Init(thor.Configuration, thor.Logger, thor.RegisterActions)` function. This function is called
on THOR startup and allows plugins to define the conditions when a plugin should be notified.

Plugins are notified when their conditions are matched by some analyzed data, including said data. They can then work
with this data to possibly create alerts or scan further data.

### Examples

See the `examples/` folder for several examples and ideas on how to write THOR plugins.