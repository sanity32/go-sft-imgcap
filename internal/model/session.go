package model

import "github.com/sanity32/b64img"

type Session struct {
	Sid   string         `json:"sid"`
	Descr string         `json:"descr"`
	Bb64  []b64img.Image `json:"bb64"`
}

func (ses Session) z() {
	b := ses.Bb64[0]
	b.Hash()
}
