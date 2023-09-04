package confmigrate

// migrateTo16 performs the following changes:
//
//	# BEFORE:
//	'schema_version': 15
//	'dns':
//	  # …
//	  'statistics_interval': 1
//	'statistics':
//	  # …
//	# …
//
//	# AFTER:
//	'schema_version': 16
//	'dns':
//	  # …
//	'statistics':
//	  'enabled': true
//	  'interval': 1
//	  'ignored': []
//	  # …
//	# …
//
// If statistics were disabled:
//
//	# BEFORE:
//	'schema_version': 15
//	'dns':
//	  # …
//	  'statistics_interval': 0
//	'statistics':
//	  # …
//	# …
//
//	# AFTER:
//	'schema_version': 16
//	'dns':
//	  # …
//	'statistics':
//	  'enabled': false
//	  'interval': 1
//	  'ignored': []
//	  # …
//	# …
func migrateTo16(diskConf yobj) (err error) {
	diskConf["schema_version"] = 16

	dns, ok, err := fieldVal[yobj](diskConf, "dns")
	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	stats := yobj{
		"enabled":  true,
		"interval": 1,
		"ignored":  yarr{},
	}

	const field = "statistics_interval"

	statsIvl, ok, err := fieldVal[int](dns, field)
	if err != nil {
		return err
	} else if ok {
		if statsIvl == 0 {
			// Set the interval to the default value of one day to make sure
			// that it passes the validations.
			stats["enabled"] = false
		} else {
			stats["interval"] = statsIvl
		}
		delete(dns, field)
	}

	diskConf["statistics"] = stats

	return nil
}
