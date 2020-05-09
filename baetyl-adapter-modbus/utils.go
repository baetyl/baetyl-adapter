package main

import (
	"fmt"
	"github.com/baetyl/baetyl-go/context"
)

func setDefault(cfg *Config, ctx context.Context) {
	var js []Job
	for _, job := range cfg.Jobs {
		if job.Time.Type == "" {
			job.Time.Type = "s"
		}
		if job.Time.Name == "" {
			job.Time.Name = "time"
		}
		var ms []MapConfig
		for _, m := range job.Maps {
			if job.Encoding == JsonEncoding {
				populateQuantityIfNeeds(&m)
				ms = append(ms, m)
			}
		}
		job.Maps = ms
		js = append(js, job)
	}
	cfg.Jobs = js
}

func populateQuantityIfNeeds(cfg *MapConfig) {
	switch cfg.Field.Type {
	case Bool:
		cfg.Quantity = 1
	case Int16:
		cfg.Quantity = 1
	case UInt16:
		cfg.Quantity = 1
	case Int32:
		cfg.Quantity = 2
	case UInt32:
		cfg.Quantity = 2
	case Int64:
		cfg.Quantity = 4
	case UInt64:
		cfg.Quantity = 4
	case Float32:
		cfg.Quantity = 2
	case Float64:
		cfg.Quantity = 4
	default:
	}
}

func validate(cfg Config) error {
	slaves := map[byte]struct{}{}
	for _, s := range cfg.Slaves {
		slaves[s.ID] = struct{}{}
	}
	for _, job := range cfg.Jobs {
		if _, ok := slaves[job.SlaveId]; !ok {
			return fmt.Errorf("slave id (%d) of job not defined in slave list", job.SlaveId)
		}
		for _, m := range job.Maps {
			if job.Encoding == JsonEncoding {
				if m.Field.Name == "" || m.Field.Type == "" {
					return fmt.Errorf("field name or type of map %+v shall not be empty when encoding is json", m)
				}
			} else if job.Encoding == BinaryEncoding {
				if m.Quantity == 0 {
					return fmt.Errorf("quantity of map %+v shall not be zero when encoding is binary", m)
				}
			}
		}
	}
	return nil
}
