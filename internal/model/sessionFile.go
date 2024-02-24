package model

import "github.com/sanity32/b64img"

type SessionFile struct {
	Sid    string
	Hashes []b64img.Hash
}
