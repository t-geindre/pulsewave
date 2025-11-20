package tree

import "fmt"

func formatSemiTon(v float32) string {
	return fmt.Sprintf("%.2f st", v)
}

func formatOctave(v float32) string {
	return fmt.Sprintf("%.0f oct", v/12)
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

func formatVoice(v float32) string {
	return fmt.Sprintf("%.0f voices", v)
}
