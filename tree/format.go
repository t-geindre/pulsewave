package tree

import "fmt"

func formatSemiTon(v float32) string {
	return fmt.Sprintf("%.2f st", v)
}

func formatMillisecond(v float32) string {
	return fmt.Sprintf("%.0f ms", v*1000)
}

func formatHertz(v float32) string {
	return fmt.Sprintf("%.0f Hz", v)
}

func formatLowHertz(v float32) string {
	return fmt.Sprintf("%.2f Hz", v)
}

func formatCycle(v float32) string {
	return fmt.Sprintf("%.0f%% cycle", v*100)
}

func formatCent(v float32) string {
	return fmt.Sprintf("%.1f cent", v)
}

func formatOnOff(v float32) string {
	if v < 0.5 {
		return "OFF"
	}
	return "ON"
}
