package maps

import (
	memstore "github.com/odio4u/memstore/registry/pkg/memstore"
	mapper "github.com/odio4u/memstore/registry/proto"
	"github.com/odio4u/memstore/registry/wal"
)

type RPCMap struct {
	mapper.UnimplementedMapsServer
	MemStore *memstore.MemStore
	WALer    *wal.WALer
}

var _ mapper.MapsServer = (*RPCMap)(nil)

func NewRPCMap() *RPCMap {
	return &RPCMap{
		MemStore: memstore.NewMemStore(),
	}
}
