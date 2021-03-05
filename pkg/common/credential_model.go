package common

type KlevrCredential struct {
	ID        uint64   `json:"id"`
	ZoneID    uint64   `json:"zoneId"`
	Name      string   `json:"name"`
	Value     string   `json:"value"`
	CreatedAt JSONTime `json:"createdAt"`
	UpdatedAt JSONTime `json:"updatedAt"`
}
