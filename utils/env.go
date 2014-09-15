package utils

import "os"

func GetEnv(key, def string) string {
	if ev := os.Getenv(key); len(ev) > 0 {
		return ev
	}
	return def
}
