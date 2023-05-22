package config

import (
	dbg "runtime/debug"
)

func GetVersionInfo() *dbg.BuildInfo {
	i, _ := dbg.ReadBuildInfo()
	return i
}

func GetBuildSettings(i *dbg.BuildInfo) map[string]string {
	m := make(map[string]string)
	for _, setting := range i.Settings {
		switch setting.Key {
		case "vcs.revision":
			m["revision"] = setting.Value
		case "vcs.modified":
			m["modified"] = setting.Value
		case "vcs.time":
			m["time"] = setting.Value
		}
	}
	return m
}
