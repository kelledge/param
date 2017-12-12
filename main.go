package main

import "fmt"
import "os"
import "syscall"
import "strings"
import "strconv"

var data = map[string]map[int]string{
	"DATABASE_PASSWORD": {
		1: "pa55w0rd",
		2: "l33t0",
	},
	"API_TOKEN": {
		1: "112233",
		2: "334455",
	},
}

func mapToEnv(env map[string]string) []string {
	envOs := make([]string, 10)

	for key, value := range env {
		element := strings.Join([]string{key, value}, "=")
		envOs = append(envOs, element)
	}

	return envOs
}

func envToMap(env []string) map[string]string {
	envMap := make(map[string]string)

	for _, e := range env {
		pair := strings.Split(e, "=")
		key := pair[0]
		value := strings.Join(pair[1:], "=")
		envMap[key] = value
	}

	return envMap
}

func printEnv(env map[string]string) {
	for key, value := range env {
		if strings.HasPrefix(key, "SSM_PARAM") {
			fmt.Printf("[%s]: %s\n", key, value)
		}
	}
}

func stripEnvWithPrefix(env map[string]string, prefix string) (map[string]string, map[string]string) {

	strippedEnv := make(map[string]string)
	strippedPairs := make(map[string]string)

	for key, value := range env {
		if strings.HasPrefix(key, prefix) {
			newKey := strings.TrimPrefix(key, prefix)
			strippedPairs[newKey] = value
		} else {
			strippedEnv[key] = value
		}
	}

	return strippedEnv, strippedPairs
}

func resolveLookups(env map[string]string) map[string]string {
	lookups := make(map[string]string)

	for key, value := range env {
		version, err := strconv.Atoi(value)
		if err != nil {
			return nil
		}

		if result, ok := data[key][version]; ok {
			lookups[key] = result
		}
	}

	return lookups
}

func mergeEnv(left map[string]string, right map[string]string) map[string]string {
	merged := make(map[string]string)

	for key, value := range left {
		merged[key] = value
	}

	for key, value := range right {
		merged[key] = value
	}

	return merged
}

func main() {
	fmt.Println("Hello")
	args := []string{"env"}

	env := os.Environ()
	envMap := envToMap(env)

	strippedEnv, strippedPairs := stripEnvWithPrefix(envMap, "SSM_PARAM_")
	resolvedEnv := resolveLookups(strippedPairs)

	mergedEnv := mergeEnv(strippedEnv, resolvedEnv)

	syscall.Exec("/usr/bin/env", args, mapToEnv(mergedEnv))
}
