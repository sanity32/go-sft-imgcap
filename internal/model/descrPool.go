package model

import (
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/gofrs/uuid"
	"github.com/sanity32/b64img"
)

func NewNormalPool(descr Descr) DescrPool {
	return DescrPool{
		Descr:     descr,
		Capacity:  200,
		Threshold: 50,
	}
}

type DescrPool struct {
	Descr     Descr
	Capacity  int
	Threshold int
}

func (d DescrPool) EnoughSamples() bool {
	return len(d.List()) >= d.Capacity
}

func (d DescrPool) List() []os.DirEntry {
	dir := d.Descr.Folder()
	if !dirExist(dir) {
		return nil
	}
	if entries, err := os.ReadDir(dir); err != nil {
		return nil
	} else {
		return entries
	}
}

func (d DescrPool) ListSorted() []os.DirEntry {
	dir := d.Descr.Folder()
	if !dirExist(dir) {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	sort.Slice(entries, func(i, j int) bool {
		infoI, _ := entries[i].Info()
		infoJ, _ := entries[j].Info()
		return infoI.ModTime().Unix() > infoJ.ModTime().Unix()
	})

	return entries
}

func (dp DescrPool) RemoveOld() (remains []os.DirEntry, err error) {
	n := dp.Capacity
	l := dp.ListSorted()
	if len(l) < n {
		return l, nil
	}
	remains = l[:n]
	oldList := l[n:]
	for _, oldRec := range oldList {
		filename := path.Join(dp.Descr.Folder(), oldRec.Name())
		if err := os.Remove(filename); err != nil {
			return remains, err
		}
	}
	return
}

func (dp DescrPool) Distinct() (rr map[b64img.Hash]int, err error) {
	list, err := dp.RemoveOld()
	rr = make(map[b64img.Hash]int, len(list))
	inc := func(k b64img.Hash) {
		if _, ok := rr[k]; !ok {
			rr[k] = 0
		}
		rr[k]++
	}
	if err != nil {
		return rr, err
	}

	for _, item := range list {
		sid := item.Name()
		rec := SessionRecord{
			sid:   uuid.FromStringOrNil(sid),
			Descr: dp.Descr,
		}
		if err := rec.Load(); err != nil {
			return rr, err
		}
		for _, img := range rec.Images {
			inc(img.Hash())
		}
	}
	return
}

func (dp DescrPool) Weights(hashes []b64img.Hash) (rr []float32, err error) {
	if dp.Threshold == 0 {
		return rr, ErrThresholdIsZero
	}
	m, err := dp.Distinct()
	if err != nil {
		return rr, err
	}
	rr = make([]float32, len(hashes))
	for n, hash := range hashes {
		occurrences := m[hash]
		fmt.Printf("n: %v, hash: %v, occ: %v\n", n, hash, occurrences)
		rr[n] = float32(occurrences) / float32(dp.Threshold)
	}
	return rr, nil
}
