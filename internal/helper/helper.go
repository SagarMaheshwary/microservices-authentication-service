package helper

import (
	"path"
	"path/filepath"
	"runtime"

	"google.golang.org/grpc/metadata"
)

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))

	return filepath.Dir(d)
}

func GetFromMetadata(md metadata.MD, k string) (string, bool) {
	v := md.Get(k)

	if len(v) == 0 {
		return "", false
	}

	return v[0], true
}
