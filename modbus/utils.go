package main

import (
	"errors"
	"fmt"
	"reflect"
)

func setDefault(cfg *Config) {
	var js []Job
	for _, job := range cfg.Jobs {
		if job.Time.Name == "" {
			job.Time.Name = SysTime
		}
		if job.Time.Type == "" {
			job.Time.Type = IntegerTime
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

func validateJobs(v interface{}, param string) error {
	if reflect.ValueOf(v).Kind() == reflect.Slice {
		jobs, ok := v.([]Job)
		if !ok {
			return errors.New("only support job array")
		}
		for _, job := range jobs {
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
				if _, ok := SysName[m.Field.Name]; ok {
					return fmt.Errorf("can not define sys field name: %s", m.Field.Name)
				}
				if _, ok := SysType[m.Field.Type]; !ok {
					return fmt.Errorf("unsupported field type: %s", m.Field.Type)
				}
			}
		}
	} else {
		return errors.New("unsupported type")
	}
	return nil
}
