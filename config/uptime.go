package config


func FormatUptime(d time.Duration) string {
	seconds := int(d.Seconds()) % 60
	minutes := int(d.Minutes()) % 60
	hours := int(d.Hours()) % 24
	days := int(d.Hours()) / 24

	var result string
	if days > 0 {
		result += fmt.Sprintf("%d Day", days)
		if days > 1 {
			result += "s"
		}
		result += " "

		if hours > 0 {
			result += fmt.Sprintf("%d Hour", hours)
			if hours > 1 {
				result += "s"
			}
			result += " "
		}
		if minutes > 0 {
			result += fmt.Sprintf("%d Minute", minutes)
			if minutes > 1 {
				result += "s"
			}
		}
	} else {
		if hours > 0 {
			result += fmt.Sprintf("%d Hour", hours)
			if hours > 1 {
				result += "s"
			}
			result += " "
		}
		if minutes > 0 {
			result += fmt.Sprintf("%d Minute", minutes)
			if minutes > 1 {
				result += "s"
			}
			result += " "
		}
		if seconds > 0 || result == "" {
			result += fmt.Sprintf("%d Second", seconds)
			if seconds != 1 {
				result += "s"
			}
		}
	}

	return result
}