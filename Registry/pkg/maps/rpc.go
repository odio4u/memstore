package maps

import (
	memstore "github.com/Purple-House/memstore/registry/pkg/memstore"
	mapper "github.com/Purple-House/memstore/registry/proto"
)

type RPCMap struct {
	mapper.UnimplementedMapsServer
	MemStore *memstore.MemStore
}

var _ mapper.MapsServer = (*RPCMap)(nil)

func NewRPCMap() *RPCMap {
	return &RPCMap{
		MemStore: memstore.NewMemStore(),
	}
}
