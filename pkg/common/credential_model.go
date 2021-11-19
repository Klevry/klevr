package common

import "github.com/Klevry/klevr/pkg/serialize"

type KlevrCredential struct {
	ID        uint64             `json:"id"`
	ZoneID    uint64             `json:"zoneId"`
	Key       string             `json:"key"`
	Value     string             `json:"value"`
	Hash      string             `json:"hash"`
	CreatedAt serialize.JSONTime `json:"createdAt"`
	UpdatedAt serialize.JSONTime `json:"updatedAt"`
}
