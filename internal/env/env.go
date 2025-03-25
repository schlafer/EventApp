package env

import (
	"os"
	"strconv"
)

func GetEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

/*
The GetEnvString and GetEnvInt functions are used
to get the value of an environment variable.
If the environment variable is not set,
the function returns the default value.
*/
