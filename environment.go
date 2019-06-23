package goenv

import (
	"strings"
)

type environment map[string]string

func parseEnvironment(prefix string, environFunc func() []string) environment {
	environ := environFunc()
	output := make(environment, len(environ))
	for _, keyvalue := range environ {
		parts := strings.SplitN(keyvalue, "=", 2)
		name := parts[0]
		value := parts[1]
		if prefix != "" {
			if !strings.HasPrefix(name, prefix) {
				continue
			}
			name = strings.TrimPrefix(name, prefix)
		}
		output[name] = value
	}
	return output
}

func (e environment) LookupEnv(name string) (string, bool) {
	value, ok := e[name]
	return value, ok
}
