package sysenv

import "os"

func GetEnvValue(name, defaultValue string) string {
	if val, ok := os.LookupEnv(name); !ok {
		return defaultValue
	} else {
		return val
	}
}
