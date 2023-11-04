package logger

import "os"

type Config struct {
	LogLevel   string
	OutputPath string
	ErrorPath  string
}

func (c *Config) ParseConfigFromEnv() {
	parseEnv := func(envName, defValue string) string {
		if v := os.Getenv(envName); v != "" {
			return v
		}
		return defValue
	}

	c.LogLevel = parseEnv("LOG_LEVEL", "info")
	c.OutputPath = parseEnv("LOG_OUTPUT_PATH", "stdout")
	c.ErrorPath = parseEnv("LOG_OUTPUT_ERROR_PATH", "stderr")
}
