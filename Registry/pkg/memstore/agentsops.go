package memstore

import "fmt"

func (mem *MemStore) AddAgent(region string, agent *AgentData) (*AgentData, *GatewayData, error) {
	data := mem.RegionExist(region)

	data.Mu.Lock()
	defer data.Mu.Unlock()

	gateway, exist := mem.GetGateway(region, agent.GatewayAddress)
	if !exist {
		return &AgentData{}, nil, fmt.Errorf("gateway %s not found in region %s", agent.GatewayAddress, region)
	}

	_, exist = data.Agents[agent.AgentDomain]
	if exist {
		fmt.Printf("Agent %s already exists in region %s\n", agent.AgentDomain, region)
		agent_data := data.Agents[agent.AgentDomain]
		agent_data.GatewayID = gateway.GatewayID
		agent_data.GatewayIP = gateway.GatewayIP
		agent_data.GatewayAddress = gateway.GatewayAddress
		return agent_data, gateway, nil
	}

	data.Agents[agent.AgentDomain] = agent
	return agent, gateway, nil
}
