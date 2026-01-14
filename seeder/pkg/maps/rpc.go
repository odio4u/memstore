package maps

import (
	mapper "github.com/odio4u/agni-schema/maps"
	memstore "github.com/odio4u/memstore/seeder/pkg/memstore"
	"github.com/odio4u/memstore/seeder/wal"
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
