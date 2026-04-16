package main

import (
	"math/rand"
)

func randUse() string {
	if len(PluginsGo) == 0 {
		return ""
	}
	pluginUses := PluginsGo[rand.Intn(len(PluginsGo))].Usage
	if pluginUses == nil || len(pluginUses) == 0 {
		return ""
	}
	return pluginUses[rand.Intn(len(pluginUses))]
}
