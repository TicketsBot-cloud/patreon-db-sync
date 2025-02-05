package model

import "time"

type (
	Tier     int32
	SkuLabel string

	Entitlement struct {
		Tier          Tier      `json:"tier"`
		Label         SkuLabel  `json:"label"`
		PatreonTierId uint64    `json:"patreon_tier_id"`
		IsLegacy      bool      `json:"is_legacy"`
		Priority      int32     `json:"priority"`
		ExpiresAt     time.Time `json:"expires_at"`
	}
)
