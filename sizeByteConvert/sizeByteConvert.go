package sizeByteConvert

import (
	"github.com/fatih/color"
)

const (
	_       = iota
	KB float64 = 1 << (10 * iota)
	MB
	GB
	TB
)
func Convert(value int64) string {
	b := float64(value)
	switch {
	case b >= TB:
		return color.RedString("%.2fTB", b/TB)
	case b >= GB:
		return color.MagentaString("%.2fGB", b/GB)
	case b >= MB:
		return color.GreenString("%.2fMB", b/MB)
	case b >= KB:
		return color.CyanString("%.2fKB", b/KB)
	}
	return color.BlueString("%.2fB", b)
}
