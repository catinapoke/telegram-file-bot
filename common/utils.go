package common

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func GetEnvFromFile(path_env string) (string, error) {
	path := os.Getenv(path_env)
	data, err := os.ReadFile(path)

	if err != nil {
		return "", fmt.Errorf("can't retrieve env value %s - '%s': %w", path, path_env, err)
	}

	return strings.TrimRight(string(data), "\n"), nil
}

func GetEnv(name string) (string, error) {
	data := os.Getenv(name)
	if data == "" {
		return "", fmt.Errorf("empty enviroment variable %s", name)
	}

	return data, nil
}

func GenerateRandomString(length int) string {
	data := "abcdefghijklmnopqrstuvwxyz01234567890$#@&/*-+"
	builder := strings.Builder{}
	builder.Grow(length)
	for i := 0; i < length; i++ {
		builder.WriteByte(data[rand.Intn(len(data))])
	}

	return builder.String()
}
