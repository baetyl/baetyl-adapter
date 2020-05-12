package main

func setDefault(cfg *Config) {
	var js []Job
	for _, job := range cfg.Jobs {
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
