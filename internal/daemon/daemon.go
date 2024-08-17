package daemon

import (
	"context"
	"fmt"
	"github.com/TicketsBot/database"
	"github.com/TicketsBot/patreon-db-sync/internal/config"
	"github.com/TicketsBot/patreon-db-sync/internal/patreonproxy"
	"github.com/TicketsBot/patreon-db-sync/pkg/model"
	"go.uber.org/zap"
	"time"
)

type Daemon struct {
	config  config.Config
	db      *database.Database
	logger  *zap.Logger
	patreon *patreonproxy.Client
}

const ExecutionTimeout = time.Second * 90

func NewDaemon(config config.Config, db *database.Database, logger *zap.Logger, patreon *patreonproxy.Client) *Daemon {
	return &Daemon{
		config:  config,
		db:      db,
		logger:  logger,
		patreon: patreon,
	}
}

func (d *Daemon) Start() error {
	d.logger.Info("Starting daemon", zap.Duration("frequency", d.config.RunFrequency))
	ctx := context.Background()

	timer := time.NewTimer(d.config.RunFrequency)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			start := time.Now()
			if err := d.doRun(ctx, ExecutionTimeout); err != nil {
				d.logger.Error("Failed to run", zap.Error(err))
			}

			d.logger.Info("Run completed", zap.Duration("duration", time.Since(start)))

			timer.Reset(d.config.RunFrequency)
		case <-ctx.Done():
			d.logger.Info("Shutting down daemon")
			return nil
		}
	}
}

func (d *Daemon) doRun(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return d.RunOnce(ctx)
}

func (d *Daemon) RunOnce(ctx context.Context) error {
	d.logger.Debug("Running synchronisation")

	d.logger.Debug("Fetching entitlements")
	res, err := d.patreon.ListEntitlements(ctx, true)
	if err != nil {
		return err
	}

	allowRemovals := true
	if len(res.Entitlements) < d.config.MinEntitlementsThreshold {
		d.logger.Warn("Number of entitlements below threshold", zap.Int("count", len(res.Entitlements)))
		allowRemovals = false
	}

	if res.LastPollTime.Before(time.Now().Add(-time.Hour)) {
		d.logger.Warn("Last poll time is older than 1 hour", zap.Time("last_poll_time", res.LastPollTime))
		allowRemovals = false
	}

	if !allowRemovals {
		d.logger.Warn("Continuing, but not removing entitlements")
	}

	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	for userId, entitlements := range res.Entitlements {
		if len(entitlements) == 0 {
			d.logger.Warn("User has no entitlements", zap.Uint64("user_id", userId))
			continue
		}

		topEntitlement := d.findTopEntitlement(entitlements)
		if topEntitlement.ExpiresAt.Add(time.Hour * 24 * time.Duration(d.config.GracePeriodDays)).Before(time.Now()) {
			d.logger.Debug("Received expired entitlement", zap.Uint64("user_id", userId), zap.Any("entitlement", topEntitlement))
			continue
		}

		d.logger.Debug("Updating entitlement", zap.Uint64("user_id", userId), zap.Any("entitlement", topEntitlement))

		if err := d.db.LegacyPremiumEntitlements.SetEntitlement(ctx, tx, database.LegacyPremiumEntitlement{
			UserId:    userId,
			TierId:    int32(topEntitlement.Tier),
			SkuLabel:  string(topEntitlement.Label),
			ExpiresAt: topEntitlement.ExpiresAt,
		}); err != nil {
			d.logger.Error("Failed to set entitlement", zap.Uint64("user_id", userId), zap.Error(err))
			return err
		}
	}

	d.logger.Info("Updated entitlements", zap.Int("count", len(res.Entitlements)))

	if allowRemovals {
		allEntitlements, err := d.db.LegacyPremiumEntitlements.ListAll(ctx, tx)
		if err != nil {
			return err
		}

		var removedCount int
		for _, existingEntitlement := range allEntitlements {
			var valid bool

			userEntitlements, ok := res.Entitlements[existingEntitlement.UserId]
			if ok {
				for _, entitlement := range userEntitlements {
					// Match entitlement: tier should match, as we've already run the update
					if entitlement.Label == model.SkuLabel(existingEntitlement.SkuLabel) &&
						entitlement.ExpiresAt.Add(time.Hour*24*time.Duration(d.config.GracePeriodDays)).After(time.Now()) {
						valid = true
						break
					}
				}
			}

			if !valid {
				d.logger.Debug("Removing entitlement", zap.Uint64("user_id", existingEntitlement.UserId))
				if err := d.db.LegacyPremiumEntitlements.Delete(ctx, tx, existingEntitlement.UserId, existingEntitlement.SkuLabel); err != nil {
					d.logger.Error("Failed to remove entitlement", zap.Uint64("user_id", existingEntitlement.UserId), zap.Error(err))
					return err
				}

				d.logger.Info("Removed entitlement", zap.Uint64("user_id", existingEntitlement.UserId), zap.Time("expires_at", existingEntitlement.ExpiresAt))

				removedCount++
			}
		}

		if removedCount > d.config.MaxRemovalsThreshold {
			d.logger.Error("Too many entitlements flagged for removal", zap.Int("count", removedCount))
			return fmt.Errorf("too many entitlements flagged for removal: %d", removedCount)
		}

		d.logger.Info("Removed entitlements", zap.Int("count", removedCount))
	}

	return tx.Commit(ctx)
}

func (d *Daemon) findTopEntitlement(entitlements []model.Entitlement) model.Entitlement {
	top := entitlements[0]

	for _, entitlement := range entitlements {
		if entitlement.Tier >= top.Tier {
			top = entitlement
		}
	}

	return top
}
