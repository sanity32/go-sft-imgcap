package model

import (
	"os"
	"path"

	"github.com/sanity32/b64img"
)

var ASSETS_DIR = func() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}
	return path.Join(homeDir, "sft-assets")
}()

var HashDir = b64img.HashDir(path.Join(ASSETS_DIR, "images"))

func dirExist(dir string) bool {
	stat, err := os.Stat(dir)
	return err == nil && stat.IsDir()
}
