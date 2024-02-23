package model

import (
	"os"
	"path"

	"github.com/sanity32/b64img"
)

func SaveImg(img b64img.Image) error {
	dir := path.Join("assets", "images")
	os.MkdirAll(dir, 0644)
	filename := path.Join(dir, string(img.Hash())+".jpg")
	return img.SaveJpeg(filename)
}

func LoadImg(hash b64img.Hash) {
	dir := path.Join("assets", "images")
	filename := path.Join(dir, string(hash)+".jpg")
	b64img.Load()
}
