package patreonproxy

import (
	"github.com/TicketsBot/patreon-db-sync/pkg/model"
	"time"
)

type ListEntitlementsResponse struct {
	Entitlements map[uint64][]model.Entitlement `json:"entitlements"`
	LastPollTime time.Time                      `json:"last_poll_time"`
}
