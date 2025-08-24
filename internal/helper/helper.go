package helper

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"google.golang.org/grpc/metadata"
)

func GetRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))

	return filepath.Dir(d)
}

func GetGRPCMetadataValue(md metadata.MD, k string) (string, bool) {
	v := md.Get(k)
	if len(v) == 0 {
		return "", false
	}

	return v[0], true
}

func GetEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

func GetEnvInt(key string, defaultVal int) int {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return val
	}

	return defaultVal
}

func GetEnvDurationSeconds(key string, defaultVal time.Duration) time.Duration {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return time.Duration(val) * time.Second
	}

	return defaultVal * time.Second
}
