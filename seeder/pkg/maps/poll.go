package maps

import (
	"context"
	"fmt"

	mapper "github.com/odio4u/agni-schema/maps"
)

func (rpc *RPCMap) ResolveGatewayForAgent(ctx context.Context, req *mapper.GatewayHandshake) (*mapper.MultipleGateways, error) {

	gateways := rpc.MemStore.GetTopKGateways("global", 10)

	var gatewayResponses []*mapper.GatewayResponse
	for _, gateway := range gateways {
		fmt.Println("gateway details:", gateway.GatewayID, gateway.Wssport)
		gatewayResponses = append(gatewayResponses, &mapper.GatewayResponse{
			GatewayId:      gateway.GatewayID,
			GatewayIp:      gateway.GatewayIP,
			GatewayAddress: gateway.GatewayAddress,
			GatewayPort:    gateway.GatewayPort,
			WssPort:        gateway.Wssport,
			Identity:       gateway.VerifiableHash,
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

func (rpc *RPCMap) ResolveGatewayForProxy(ctx context.Context, req *mapper.ProxyMapping) (*mapper.AgentResponse, error) {

	agent, exist := rpc.MemStore.GetAgent(
		req.AgentDomain,
		req.Region,
	)

	fmt.Println("ResolveGatewayForProxy capacity:")

	if exist {
		return &mapper.AgentResponse{
			AgentId:        agent.AgentID,
			AgentDomain:    agent.AgentDomain,
			GatewayId:      agent.GatewayID,
			GatewayAddress: agent.GatewayAddress,
			GatewayIp:      agent.GatewayIP,
			GatewayPort:    agent.GatewayPort,
			WssPort:        agent.Wssport,
			Identity:       agent.VerifiableHash,
		}, nil
	}
	return &mapper.AgentResponse{
		Error: &mapper.Error{
			Code:    2,
			Message: "gateway not found",
		},
	}, nil
}
