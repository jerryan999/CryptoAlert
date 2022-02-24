package utils

func GetAlertQueueKey(crypto string, direction bool) string {
	if direction {
		return "alerts:" + crypto + ":" + "gt"
	}
	return "alerts:" + crypto + ":" + "lt"
}
