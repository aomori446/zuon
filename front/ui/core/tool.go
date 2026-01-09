package core

import "fmt"

func FormatBytes(s int) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB"}
	i := 0
	v := float64(s)
	for v > 1024 && i < len(sizes)-1 {
		v /= 1024
		i++
	}
	return fmt.Sprintf("%.2f %s", v, sizes[i])
}
