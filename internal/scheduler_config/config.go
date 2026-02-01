package scheduler_config

import (
	"fmt"
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

