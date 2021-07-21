package common

type KlevrCredential struct {
	ID        uint64   `json:"id"`
	ZoneID    uint64   `json:"zoneId"`
	Key       string   `json:"key"`
	Value     string   `json:"value"`
	Hash      string   `json:"hash"`
	CreatedAt JSONTime `json:"createdAt"`
	UpdatedAt JSONTime `json:"updatedAt"`
}
