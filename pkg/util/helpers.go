package util

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// ParseColor 解析颜色字符串（如"#FF0000"）为RGBA
func ParseColor(hexColor string) (color.RGBA, error) {
	if !strings.HasPrefix(hexColor, "#") || len(hexColor) != 7 {
		return color.RGBA{}, fmt.Errorf("无效的颜色格式: %s", hexColor)
	}

	r, err := strconv.ParseUint(hexColor[1:3], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}

	g, err := strconv.ParseUint(hexColor[3:5], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}

	b, err := strconv.ParseUint(hexColor[5:7], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	}, nil
}

// FormatTime 将时间（秒）格式化为字符串（例如 "00:01:30.500"）
func FormatTime(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := (int(seconds) % 3600) / 60
	secs := int(seconds) % 60
	millisecs := int((seconds - float64(int(seconds))) * 1000)

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, secs, millisecs)
}

// Clamp 将值限制在指定范围内
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// LinearMap 将值从一个范围线性映射到另一个范围
func LinearMap(value, fromMin, fromMax, toMin, toMax float64) float64 {
	if fromMax == fromMin {
		return toMin
	}

	// 计算映射比例
	ratio := (value - fromMin) / (fromMax - fromMin)
	return toMin + ratio*(toMax-toMin)
}
