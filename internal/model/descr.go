package model

import (
	"os"
	"path"
	"sort"

	"github.com/gofrs/uuid"
	"github.com/sanity32/b64img"
)

type Descr string

func (d Descr) String() string {
	return string(d)
}

func (descr Descr) Folder() string {
	return path.Join(ASSETS_DIR, "descr", string(descr))
}

func (descr Descr) ListSids() (rr []string) {
	dir := descr.Folder()
	if !dirExist(dir) {
		return
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		infoI, _ := entries[i].Info()
		infoJ, _ := entries[j].Info()
		return infoI.ModTime().Unix() > infoJ.ModTime().Unix()
	})

	for _, e := range entries {
		if !e.IsDir() {
			rr = append(rr, e.Name())
		}
	}
	return
}

func (descr Descr) Read() (rr []SessionRecord, err error) {
	for _, sid := range descr.ListSids() {
		rec := SessionRecord{
			sid:   uuid.FromStringOrNil(sid),
			Descr: descr,
		}
		if err := rec.Load(); err != nil {
			return rr, err
		}
		rr = append(rr, rec)
	}
	return
}

func (descr Descr) Distinct() map[b64img.Image]int {
	var rr = map[b64img.Image]int{}
	for _, sid := range descr.ListSids() {
		rec := SessionRecord{
			sid:   uuid.FromStringOrNil(sid),
			Descr: descr,
		}
		if err := rec.Load(); err == nil {
			for _, img := range rec.Images {
				if _, ok := rr[img]; !ok {
					rr[img] = 0
				}
				rr[img]++
			}
		}
	}
	return rr
}
