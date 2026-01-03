package maps

import (
	"context"

	memstore "github.com/Purple-House/memstore/registry/pkg/memstore"
	mapper "github.com/Purple-House/memstore/registry/proto"
	walpb "github.com/Purple-House/memstore/registry/wal/proto"

	"github.com/google/uuid"
)

func (rpc *RPCMap) RegisterGateway(ctx context.Context, req *mapper.GatewayPutRequest) (*mapper.GatewayResponse, error) {

	if req.GatewayIp == "" || req.GatewayPort == 0 || req.VerifiableCredHash == "" {
		return &mapper.GatewayResponse{
			Error: &mapper.Error{
				Code:    1,
				Message: "invalid gateway registration request",
			},
		}, nil
	}

	gatewayData := &memstore.GatewayData{
		GatewayIP:      req.GatewayIp,
		GatewayID:      uuid.New().String(),
		GatewayPort:    req.GatewayPort,
		VerifiableHash: req.VerifiableCredHash,
		Capacity: memstore.Capacity{
			CPU:     req.Capacity.Cpu,
			Memory:  req.Capacity.Memory,
			Storage: req.Capacity.Storage,
		},
	}

	region := req.Region
	if region == "" {
		region = "global"
	}

	data, err := rpc.MemStore.AddGateway(
		region,
		gatewayData,
	)
	if err != nil {
		return &mapper.GatewayResponse{
			Error: &mapper.Error{
				Code:    1,
				Message: err.Error(),
			},
		}, nil
	}

	// this should be zero lock write to WAL
	err = rpc.WALer.Append(&walpb.WalRecord{

		Op: walpb.Operation_OP_PUT_GATEWAY,
		Gateway: &walpb.GatewayPutRequest{
			Region:             region,
			GatewayIp:          data.GatewayIP,
			GatewayId:          data.GatewayID,
			GatewayPort:        data.GatewayPort,
			GatewayAddress:     data.GatewayAddress,
			VerifiableCredHash: data.VerifiableHash,
			Capacity: &walpb.Capacity{
				Cpu:     data.Capacity.CPU,
				Memory:  data.Capacity.Memory,
				Storage: data.Capacity.Storage,
			},
		},
	})

	if err != nil {
		return &mapper.GatewayResponse{
			Error: &mapper.Error{
				Code:    1,
				Message: err.Error(),
			},
		}, nil
	}
	return &mapper.GatewayResponse{
		GatewayId:      data.GatewayID,
		GatewayIp:      data.GatewayIP,
		GatewayAddress: data.GatewayAddress,
	}, nil
}

func (rpc *RPCMap) ResolveGatewayForProxy(ctx context.Context, req *mapper.GatewayProxy) (*mapper.GatewayResponse, error) {

	gateway, exist := rpc.MemStore.GetGateway(
		req.AgentDomain,
		req.Region,
	)

	if exist {
		return &mapper.GatewayResponse{
			GatewayId:      gateway.GatewayID,
			GatewayIp:      gateway.GatewayIP,
			GatewayAddress: gateway.GatewayAddress,
		}, nil
	}
	return &mapper.GatewayResponse{
		Error: &mapper.Error{
			Code:    2,
			Message: "gateway not found",
		},
	}, nil

}
