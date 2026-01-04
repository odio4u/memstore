package maps

import (
	"context"

	memstore "github.com/Purple-House/memstore/registry/pkg/memstore"
	mapper "github.com/Purple-House/memstore/registry/proto"
	"github.com/google/uuid"
)

func (rpc *RPCMap) RegisterAgent(ctx context.Context, req *mapper.AgentConnectionRequest) (*mapper.AgentResponse, error) {

	agentData := &memstore.AgentData{
		AgentDomain: req.AgentDomain,
		AgentID:     uuid.New().String(),
	}

	agent, _, err := rpc.MemStore.AddAgent(req.Region, agentData)
	if err != nil {
		return &mapper.AgentResponse{
			Error: &mapper.Error{
				Code:    1,
				Message: err.Error(),
			},
		}, nil
	}
	return &mapper.AgentResponse{
		AgentId:     agent.AgentID,
		AgentDomain: agent.AgentDomain,
		Error:       nil,
	}, nil
}

func (rpc *RPCMap) ConnectAgentTogateway(ctx context.Context, req *mapper.AgentConnect) (*mapper.AgentResponse, error) {

	agentData := &memstore.AgentData{
		AgentDomain: req.AgentDomain,
		GatewayID:   req.GatewayId,
	}

	agent, _, err := rpc.MemStore.AddAgent("global", agentData)
	if err != nil {
		return &mapper.AgentResponse{
			Error: &mapper.Error{
				Code:    1,
				Message: err.Error(),
			},
		}, nil
	}
	return &mapper.AgentResponse{
		AgentId:     agent.AgentID,
		AgentDomain: agent.AgentDomain,
		Error:       nil,
	}, nil
}

func (rpc *RPCMap) ResolveGatewayForAgent(ctx context.Context, req *mapper.GatewayHandshake) (*mapper.MultipleGateways, error) {

	gateways := rpc.MemStore.GetTopKGateways("global", 10)

	var gatewayResponses []*mapper.GatewayResponse
	for _, gateway := range gateways {
		gatewayResponses = append(gatewayResponses, &mapper.GatewayResponse{
			GatewayId:      gateway.GatewayID,
			GatewayIp:      gateway.GatewayIP,
			GatewayAddress: gateway.GatewayAddress,
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
