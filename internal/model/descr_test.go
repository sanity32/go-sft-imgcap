package model

import (
	"fmt"
	"os"
	"testing"

	"github.com/sanity32/b64img"
)

func TestDescr_ListSids(t *testing.T) {
	descr := Descr("Отметьте все изображения с шапками")

	pool := NewNormalPool(descr)
	m, err := pool.Distinct()
	if err != nil {
		t.Fatal(err.Error())
	}
	os.Mkdir(descr.String(), 0644)

	for hash, n := range m {
		im, err := HashDir.Read(hash)
		if err != nil {
			t.Fatal(err.Error())
		}
		f := fmt.Sprintf("%v/%03d_%v.jpg", descr, n, im.Hash())
		im.SaveJpeg(f)
	}
}

func TestWeights(t *testing.T) {
	descr := Descr("Отметьте все изображения с шапками")
	pool := NewNormalPool(descr)
	testHashes := []b64img.Hash{
		"8feb88da44b38386cbc95a97406dddad",
		"83e5df5a02f8654b5b86031a50acf122",
		"5be328c24d8d6ee9498a83cb9cf10154",
		"3d15b9b0f2a9a61b452899b101ce0f10",
	}
	ww, err := pool.Weights(testHashes)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(ww)
}
