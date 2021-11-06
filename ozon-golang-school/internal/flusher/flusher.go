package flusher

import (
	"github.com/onsi/ginkgo"
	"github.com/ozonva/ova-link-api/internal/link"
	"github.com/ozonva/ova-link-api/internal/repo"
	"github.com/ozonva/ova-link-api/internal/utils"
)

type Flusher interface {
	Flush(entities []link.Link) []link.Link
}

type flusher struct {
	chunkSize  uint
	entityRepo repo.Repo
}

func NewFlusher(chunkSize uint, entityRepo repo.Repo) Flusher {
	return &flusher{
		chunkSize:  chunkSize,
		entityRepo: entityRepo,
	}
}

func (f *flusher) Flush(entities []link.Link) []link.Link {
	defer ginkgo.GinkgoRecover()
	unprocessedEntities := make([]link.Link, 0, len(entities))
	for _, batch := range utils.SliceChunkLink(entities, f.chunkSize) {
		err := f.entityRepo.AddEntities(batch)
		if err != nil {
			unprocessedEntities = append(unprocessedEntities, batch...)
		}
	}

	if len(unprocessedEntities) > 0 {
		return unprocessedEntities
	}

	return nil
}
