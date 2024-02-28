package model

import (
	"os"
	"path"

	"github.com/sanity32/b64img"
)

func SaveImg(img b64img.Image) error {
	dir := path.Join(ASSETS_DIR, "images")
	os.MkdirAll(dir, 0777)
	filename := path.Join(dir, string(img.Hash())+".jpg")
	return img.SaveJpeg(filename)
}

func LoadImg(hash b64img.Hash) {
	dir := path.Join(ASSETS_DIR, "images")
	filename := path.Join(dir, string(hash)+".jpg")
	b64img.Load(filename)
}
