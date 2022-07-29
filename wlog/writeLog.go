package wlog

import (
	"io/ioutil"
	"path"
)

func WriteLog(data []byte, outPath string) error {
	return ioutil.WriteFile(path.Join(outPath, "log.txt"), data, 0644)
}
