package jobs

import (
	"os"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type CronJob struct {
	JobConfigs []JobConfig `yaml:"jobs"`
}
type JobConfig struct {
	Schedule string `yaml:"schedule"`
	Func     string `yaml:"func"`
}

func LoadConfig(path string) (*CronJob, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config CronJob
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func InitJobs(
	c *cron.Cron,
	cfg *CronJob,
	db *gorm.DB,
) error {
	var err error
	for _, job := range cfg.JobConfigs {
		switch job.Func {
		case "abc":
			_, err = c.AddFunc(job.Schedule, func() {
				Job1(db)
			})
		default:
			continue
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func Job1(db *gorm.DB) {}
