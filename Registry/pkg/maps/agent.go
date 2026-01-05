package maps

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	memstore "github.com/Purple-House/memstore/registry/pkg/memstore"
	mapper "github.com/Purple-House/memstore/registry/proto"
	walpb "github.com/Purple-House/memstore/registry/wal/proto"
)

func (rpc *RPCMap) RegisterAgent(ctx context.Context, req *mapper.AgentConnectionRequest) (*mapper.AgentResponse, error) {

	if req.VerifiableCredHash == "" || req.AgentDomain == "" || req.GatewayId == "" || req.Region == "" {
		return &mapper.AgentResponse{
			Error: &mapper.Error{
				Code:    1,
				Message: "invalid agent registration request",
			},
		}, nil
	}

	identityBytes := sha256.Sum256([]byte(
		req.VerifiableCredHash + "|" + req.AgentDomain,
	))

	identity := hex.EncodeToString(identityBytes[:])
	fmt.Println("identity:", identity)

	agentData := &memstore.AgentData{
		AgentDomain:    req.AgentDomain,
		AgentID:        identity,
		GatewayID:      req.GatewayId,
		VerifiableHash: req.VerifiableCredHash,
	}

	agent, gateway, err := rpc.MemStore.AddAgent(req.Region, agentData)
	if err != nil {
		return &mapper.AgentResponse{
			Error: &mapper.Error{
				Code:    1,
				Message: err.Error(),
			},
		}, nil
	}

	err = rpc.WALer.Append(&walpb.WalRecord{
		Op: walpb.Operation_OP_PUT_AGENT,
		Agent: &walpb.AgentConnectionRequest{
			VerifiableCredHash: agent.VerifiableHash,
			AgentDomain:        agent.AgentDomain,
			GatewayId:          agent.GatewayID,
			Region:             req.Region,
			GatewayAddress:     gateway.GatewayAddress,
			AgentId:            agent.AgentID,
		},
	})

	return &mapper.AgentResponse{
		AgentId:        agent.AgentID,
		AgentDomain:    agent.AgentDomain,
		GatewayId:      agent.GatewayID,
		GatewayAddress: agent.GatewayAddress,
		GatewayIp:      agent.GatewayIP,
		GatewayPort:    agent.GatewayPort,
		WssPort:        agent.Wssport,
		Capacity: &mapper.Capacity{
			Cpu:     gateway.Capacity.CPU,
			Memory:  gateway.Capacity.Memory,
			Storage: gateway.Capacity.Storage,
		},
		Error: nil,
	}, nil
}
