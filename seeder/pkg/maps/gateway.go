package maps

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	walpb "github.com/odio4u/agni-schema/wal"
	memstore "github.com/odio4u/memstore/seeder/pkg/memstore"
	mapper "github.com/odio4u/memstore/seeder/proto"
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

	identityBytes := sha256.Sum256([]byte(
		req.VerifiableCredHash + "|" + req.GatewayIp,
	))

	identity := hex.EncodeToString(identityBytes[:])

	gatewayData := &memstore.GatewayData{
		GatewayIP:      req.GatewayIp,
		GatewayID:      identity,
		GatewayPort:    req.GatewayPort,
		VerifiableHash: req.VerifiableCredHash,
		Wssport:        req.WssPort,
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
			WssPort:            data.Wssport,
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
		GatewayPort:    data.GatewayPort,
		WssPort:        data.Wssport,
		Identity:       data.VerifiableHash,
		Capacity: &mapper.Capacity{
			Cpu:     data.Capacity.CPU,
			Memory:  data.Capacity.Memory,
			Storage: data.Capacity.Storage,
		},
		Error: nil,
	}, nil
}
