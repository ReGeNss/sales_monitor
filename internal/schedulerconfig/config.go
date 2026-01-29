package schedulerconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Timezone   string           `yaml:"timezone"`
	Shops      []ShopConfig     `yaml:"shops"`
	Categories []CategoryConfig `yaml:"categories"`
}

type ShopConfig struct {
	ID          string `yaml:"id"`
	DefaultCron string `yaml:"default_cron"`
}

type CategoryConfig struct {
	Name            string      `yaml:"name"`
	WordsToIgnore   []string    `yaml:"words_to_ignore"`
	Differentiation [][]string  `yaml:"differentiation"`
	Jobs            []JobConfig `yaml:"jobs"`
}

type JobConfig struct {
	ID           string   `yaml:"id"`
	ShopID       string   `yaml:"shop_id"`
	URLs         []string `yaml:"urls"`
	CronOverride string   `yaml:"cron_override"`
}

type ResolvedJob struct {
	ID              string
	ShopID          string
	Cron            string
	CategoryName    string
	WordsToIgnore   []string
	Differentiation [][]string
	URLs            []string
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve config path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

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

func (c *Config) Location() (*time.Location, error) {
	if c.Timezone == "" {
		return time.Local, nil
	}
	return time.LoadLocation(c.Timezone)
}

func (c *Config) ResolveJobs() ([]ResolvedJob, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	shopCron := map[string]string{}
	for _, shop := range c.Shops {
		shopCron[shop.ID] = shop.DefaultCron
	}

	resolved := make([]ResolvedJob, 0)
	for _, category := range c.Categories {
		for _, job := range category.Jobs {
			cronExpr := job.CronOverride
			if cronExpr == "" {
				cronExpr = shopCron[job.ShopID]
			}
			resolved = append(resolved, ResolvedJob{
				ID:              job.ID,
				ShopID:          job.ShopID,
				Cron:            cronExpr,
				CategoryName:    category.Name,
				WordsToIgnore:   category.WordsToIgnore,
				Differentiation: category.Differentiation,
				URLs:            job.URLs,
			})
		}
	}

	return resolved, nil
}

func (c *Config) FindJob(jobID string) (*ResolvedJob, error) {
	if jobID == "" {
		return nil, fmt.Errorf("job id is empty")
	}
	jobs, err := c.ResolveJobs()
	if err != nil {
		return nil, err
	}
	for _, job := range jobs {
		if job.ID == jobID {
			jobCopy := job
			return &jobCopy, nil
		}
	}
	return nil, fmt.Errorf("job %q not found", jobID)
}

func validateCronExpression(expr string) error {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(expr)
	return err
}
