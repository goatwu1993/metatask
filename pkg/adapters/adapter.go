package adapters

import (
	"metatask/pkg/schema"
)

type AdaptConfig struct {
	IgnoreNotFound bool
}

type AdapaterInterface interface {
	GenerateFromMetaTaskFile(*schema.TreeRoot, *AdaptConfig) error
}
