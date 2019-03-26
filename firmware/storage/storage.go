package storage

import (
	"github.com/jdevelop/passkeeper"
)

type SeedStorageRead interface {
	LoadSeed(id string) (*passkeeper.Seed, error)
}

type SeedStorageWrite interface {
	SaveSeed(seed passkeeper.Seed) error
}

type SeedStorageList interface {
	ListSeeds() ([]string, error)
}

type SeedStorageRemove interface {
	RemoveSeed(key string) error
}
