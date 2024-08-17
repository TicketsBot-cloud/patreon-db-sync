package model

import "time"

type (
	Tier     int32
	SkuLabel string

	Entitlement struct {
		Tier      Tier      `json:"tier"`
		Label     SkuLabel  `json:"label"`
		IsLegacy  bool      `json:"is_legacy"`
		ExpiresAt time.Time `json:"expires_at"`
	}
)

const (
	TierPremium Tier = iota
	TierWhitelabel

	SkuPremium    SkuLabel = "premium"
	SkuWhitelabel SkuLabel = "whitelabel"
)
