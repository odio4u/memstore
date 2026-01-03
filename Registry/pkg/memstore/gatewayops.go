package memstore

import (
	"fmt"

	"github.com/google/btree"
)

func (mem *MemStore) AddGateway(region string, gateway *GatewayData) (GatewayData, error) {
	data := mem.RegionExist(region)

	data.Mu.Lock()
	defer data.Mu.Unlock()

	gatewayAddress := fmt.Sprintf("%s:%d", gateway.GatewayIP, gateway.GatewayPort)
	gateway.GatewayAddress = gatewayAddress

	gatewayData, exist := data.Gateways[gateway.GatewayAddress]
	if exist {
		// Remove old rank item
		oldRank := gatewayData.Capacity.Rank()
		data.ranked.Delete(&GatewayRankItem{
			Rank: oldRank,
			ID:   gateway.GatewayAddress,
		})
		// Update gateway data
		gateway.GatewayID = gatewayData.GatewayID
	}
	data.Gateways[gateway.GatewayAddress] = gateway
	data.ranked.ReplaceOrInsert(&GatewayRankItem{
		Rank: gateway.Capacity.Rank(),
		ID:   gateway.GatewayAddress,
	})

	fmt.Printf("Added gateway %s in region %s\n", gateway.GatewayAddress, region)
	return *gateway, nil
}

func (mem *MemStore) GetTopKGateways(region string, k int) []*GatewayData {
	data := mem.RegionExist(region)

	data.Mu.RLock()
	defer data.Mu.RUnlock()

	var result []*GatewayData
	count := 0
	data.ranked.Ascend(func(item btree.Item) bool {
		if count >= k {
			return false
		}
		gi := item.(*GatewayRankItem)
		result = append(result, data.Gateways[gi.ID])
		count++
		return true
	})
	return result
}

func (mem *MemStore) GetGateway(region, gatewayAddress string) (*GatewayData, bool) {
	data := mem.RegionExist(region)

	data.Mu.RLock()
	defer data.Mu.RUnlock()
	gateway, exist := data.Gateways[gatewayAddress]
	return gateway, exist
}
