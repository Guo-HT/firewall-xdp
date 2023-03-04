package utils

func ConvertProtoCode2Status(code int) (status string) {
	switch code {
	case 0:
		return "stop"
	case 1:
		return "start"
	default:
		return "start"
	}
}
