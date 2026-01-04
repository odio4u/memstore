package maps

import (
	"context"
	"fmt"

	mapper "github.com/Purple-House/memstore/registry/proto"
)

func (rpc *RPCMap) ResolveGatewayForAgent(ctx context.Context, req *mapper.GatewayHandshake) (*mapper.MultipleGateways, error) {

	gateways := rpc.MemStore.GetTopKGateways("global", 10)

	var gatewayResponses []*mapper.GatewayResponse
	for _, gateway := range gateways {
		gatewayResponses = append(gatewayResponses, &mapper.GatewayResponse{
			GatewayId:      gateway.GatewayID,
			GatewayIp:      gateway.GatewayIP,
			GatewayAddress: gateway.GatewayAddress,
			Capacity: &mapper.Capacity{
				Cpu:       gateway.Capacity.CPU,
				Memory:    gateway.Capacity.Memory,
				Storage:   gateway.Capacity.Storage,
				Bandwidth: gateway.Capacity.Bandwidth,
			},
		})
	}

	if len(gatewayResponses) == 0 {
		return &mapper.MultipleGateways{
			Gateways: []*mapper.GatewayResponse{},
			Error: &mapper.Error{
				Code:    2,
				Message: "no gateway found",
			},
		}, nil
	}
	return &mapper.MultipleGateways{
		Gateways: gatewayResponses,
		Error:    nil,
	}, nil
}

func (rpc *RPCMap) ResolveGatewayForProxy(ctx context.Context, req *mapper.GatewayProxy) (*mapper.GatewayResponse, error) {

	gateway, exist := rpc.MemStore.GetGateway(
		req.AgentDomain,
		req.Region,
	)

	fmt.Println("ResolveGatewayForProxy capacity:", gateway.Capacity.CPU, gateway.Capacity.Memory, gateway.Capacity.Storage)

	if exist {
		return &mapper.GatewayResponse{
			GatewayId:      gateway.GatewayID,
			GatewayIp:      gateway.GatewayIP,
			GatewayAddress: gateway.GatewayAddress,
			Capacity: &mapper.Capacity{
				Cpu:       gateway.Capacity.CPU,
				Memory:    gateway.Capacity.Memory,
				Storage:   gateway.Capacity.Storage,
				Bandwidth: gateway.Capacity.Bandwidth,
			},
		}, nil
	}
	return &mapper.GatewayResponse{
		Error: &mapper.Error{
			Code:    2,
			Message: "gateway not found",
		},
	}, nil

}
