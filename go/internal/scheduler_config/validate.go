package scheduler_config

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func (c *Config) Validate() error {
	if len(c.Shops) == 0 {
		return fmt.Errorf("config must include at least one shop")
	}

	if len(c.Categories) == 0 {
		return fmt.Errorf("config must include at least one category")
	}

	if _, err := c.Location(); err != nil {
		return fmt.Errorf("invalid timezone %q: %w", c.Timezone, err)
	}

	shopIDs := map[string]string{}
	for _, shop := range c.Shops {
		if shop.ID == "" {
			return fmt.Errorf("shop id is required")
		}
		if shop.DefaultCron == "" {
			return fmt.Errorf("shop %q default_cron is required", shop.ID)
		}
		if err := validateCronExpression(shop.DefaultCron); err != nil {
			return fmt.Errorf("shop %q default_cron invalid: %w", shop.ID, err)
		}
		if _, exists := shopIDs[shop.ID]; exists {
			return fmt.Errorf("duplicate shop id %q", shop.ID)
		}
		shopIDs[shop.ID] = shop.DefaultCron
	}

	jobIDs := map[string]struct{}{}
	for _, category := range c.Categories {
		if category.Name == "" {
			return fmt.Errorf("category name is required")
		}
		if len(category.Jobs) == 0 {
			return fmt.Errorf("category %q must include at least one job", category.Name)
		}
		for _, job := range category.Jobs {
			if job.ID == "" {
				return fmt.Errorf("job id is required in category %q", category.Name)
			}
			if _, exists := jobIDs[job.ID]; exists {
				return fmt.Errorf("duplicate job id %q", job.ID)
			}
			jobIDs[job.ID] = struct{}{}
			if job.ShopID == "" {
				return fmt.Errorf("job %q shop_id is required", job.ID)
			}
			if _, exists := shopIDs[job.ShopID]; !exists {
				return fmt.Errorf("job %q references unknown shop_id %q", job.ID, job.ShopID)
			}
			if len(job.URLs) == 0 {
				return fmt.Errorf("job %q must include at least one url", job.ID)
			}
			if job.CronOverride != "" {
				if err := validateCronExpression(job.CronOverride); err != nil {
					return fmt.Errorf("job %q cron_override invalid: %w", job.ID, err)
				}
			}
		}
	}

	return nil
}

func validateCronExpression(expr string) error {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(expr)
	return err
}

func (c *Config) Location() (*time.Location, error) {
	if c.Timezone == "" {
		return time.Local, nil
	}
	return time.LoadLocation(c.Timezone)
}