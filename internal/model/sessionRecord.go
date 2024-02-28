package model

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/sanity32/b64img"
)

type SessionRecord struct {
	sid    uuid.UUID
	Descr  Descr          `json:"descr"`
	Images []b64img.Image `json:"bb64"`
}

func (rec *SessionRecord) Sid() uuid.UUID {
	if rec.sid.IsNil() {
		rec.sid = uuid.Must(uuid.NewV4())
	}
	return rec.sid
}

func (rec SessionRecord) Filepath() string {
	dir := rec.Descr.Folder()
	os.MkdirAll(dir, 0777)
	return path.Join(dir, rec.Sid().String())
}

func (rec SessionRecord) Hashes() []b64img.Hash {
	rr := make([]b64img.Hash, len(rec.Images))
	for n, img := range rec.Images {
		rr[n] = img.Hash()
	}
	return rr
}

func (rec SessionRecord) Buff() (bb bytes.Buffer) {
	for n, img := range rec.Images {
		if n != 0 {
			bb.WriteString("\n")
		}
		hash := img.Hash()
		bb.WriteString(string(hash))
	}
	fmt.Println("buff:", bb.String())
	return
}

func (rec SessionRecord) Write() error {
	bb := rec.Buff()
	return os.WriteFile(rec.Filepath(), bb.Bytes(), 0644)
}

func (rec *SessionRecord) Load() error {
	bb, err := os.ReadFile(rec.Filepath())
	if err != nil {
		return err
	}
	lines := strings.Split(string(bb), "\n")
	for _, line := range lines {
		hash := b64img.Hash(line)
		img, err := MainHashDir.Read(hash)
		if err != nil {
			return err
		}
		rec.Images = append(rec.Images, img)
	}
	return nil
}

func (rec *SessionRecord) PopulateHashDir() error {
	MainHashDir.Create()
	for _, h := range rec.Images {
		if err := MainHashDir.Write(b64img.Image(h)); err != nil {
			return err
		}
	}
	return nil
}
