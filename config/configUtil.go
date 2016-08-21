package config

import (
	"os"
	"strings"
)

func getAllEnv() map[string]string {
	envs := os.Environ()
	envMap := make(map[string]string)
	for _, env := range envs {
		envKV := strings.SplitN(env, "=", 2)
		envKey, envVal := envKV[0], envKV[1]
		envMap[envKey] = envVal
	}
	return envMap
}
