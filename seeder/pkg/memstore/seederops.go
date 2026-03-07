package memstore

func (mem *MemStore) AddSeeder(seeder *SeederData) bool {

	data := mem.RegionExist(seeder.Region)

	data.Mu.Lock()
	defer data.Mu.Unlock()

	data.Seeders[seeder.SeederID] = seeder

	return true
}

func (mem *MemStore) GetSeeders(region string) []*SeederData {
	data := mem.RegionExist(region)
	result := make([]*SeederData, 0, 5)

	for _, v := range data.Seeders {
		result = append(result, v)
		if len(result) >= 5 {
			break
		}
	}

	return result
}
