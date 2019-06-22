package goenv

import (
	"strings"
)

type environment map[string]string

func parseEnvironment(environFunc func() []string) environment {
	environ := environFunc()
	output := make(environment, len(environ))
	for _, keyvalue := range environ {
		parts := strings.SplitN(keyvalue, "=", 2)
		output[parts[0]] = parts[1]
	}
	return output
}

func (e environment) LookupEnv(name string) (string, bool) {
	value, ok := e[name]
	return value, ok
}
