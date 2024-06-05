package mayo

import (
	"strconv"
	"strings"

	hotdog "github.com/neurlang/wayland/go-wayland-browser/hotdog"
)

func hexToFloatInRange(hex string) float64 {
	number, err := strconv.ParseInt(hex, 16, 0)

	if err != nil {
		panic(err)
	}

	return float64(number / 255)
}

// RGBAToColor - Transforms RGBA color string to *hotdog.ColorRGBA
// TODO - Fix this spaghetti and parse alpha values
func RGBAToColor(colorString string) *hotdog.ColorRGBA {
	var color *hotdog.ColorRGBA

	if rgbaParams.MatchString(colorString) {
		paramString := rgbaParams.FindString(colorString)
		paramString = strings.Trim(paramString, "()")

		params := strings.Split(paramString, ",")
		paramsLen := len(params)

		if paramsLen >= 3 {
			var red float64
			var green float64
			var blue float64
			var alpha float64

			if strings.HasSuffix(params[0], "%") {
				value, _ := strconv.ParseInt(strings.Trim(strings.TrimSpace(params[0]), "%"), 10, 0)
				red = float64(value / 100)
			} else if strings.Index(params[0], ".") != -1 {
				value, _ := strconv.ParseFloat(strings.TrimSpace(params[0]), 64)
				red = value
			} else {
				value, _ := strconv.Atoi(strings.TrimSpace(params[0]))
				red = float64(value / 255)
			}

			if strings.HasSuffix(params[1], "%") {
				value, _ := strconv.ParseInt(strings.Trim(strings.TrimSpace(params[1]), "%"), 10, 0)
				green = float64(value / 100)
			} else if strings.Index(params[1], ".") != -1 {
				value, _ := strconv.ParseFloat(strings.TrimSpace(params[1]), 64)
				green = value
			} else {
				value, _ := strconv.Atoi(strings.TrimSpace(params[1]))
				green = float64(value / 255)
			}

			if strings.HasSuffix(params[2], "%") {
				value, _ := strconv.ParseInt(strings.Trim(strings.TrimSpace(params[2]), "%"), 10, 0)
				blue = float64(value / 100)
			} else if strings.Index(params[2], ".") != -1 {
				value, _ := strconv.ParseFloat(strings.TrimSpace(params[2]), 64)
				blue = value
			} else {
				value, _ := strconv.Atoi(strings.TrimSpace(params[2]))
				blue = float64(value / 255)
			}

			alpha = 1

			return &hotdog.ColorRGBA{
				R: red,
				G: green,
				B: blue,
				A: alpha,
			}
		}
	}

	return color
}

// HexStringToColor - Transforms hex color string to *hotdog.ColorRGBA
func HexStringToColor(colorString string) *hotdog.ColorRGBA {
	colorString = strings.ToLower(colorString)
	colorStringLen := len(colorString)

	switch colorStringLen {
	case 9:
		return &hotdog.ColorRGBA{
			R: hexToFloatInRange(colorString[1:3]),
			G: hexToFloatInRange(colorString[3:5]),
			B: hexToFloatInRange(colorString[5:7]),
			A: hexToFloatInRange(colorString[7:9]),
		}

	case 7:
		return &hotdog.ColorRGBA{
			R: hexToFloatInRange(colorString[1:3]),
			G: hexToFloatInRange(colorString[3:5]),
			B: hexToFloatInRange(colorString[5:7]),
			A: 1,
		}

	case 5:
		return &hotdog.ColorRGBA{
			R: hexToFloatInRange(colorString[1:2] + colorString[1:2]),
			G: hexToFloatInRange(colorString[2:3] + colorString[2:3]),
			B: hexToFloatInRange(colorString[3:4] + colorString[3:4]),
			A: hexToFloatInRange(colorString[4:5] + colorString[4:5]),
		}

	case 4:
		return &hotdog.ColorRGBA{
			R: hexToFloatInRange(colorString[1:2] + colorString[1:2]),
			G: hexToFloatInRange(colorString[2:3] + colorString[2:3]),
			B: hexToFloatInRange(colorString[3:4] + colorString[3:4]),
			A: 1,
		}

	default:
		return &hotdog.ColorRGBA{}
	}
}

// MapCSSColor - Transforms css color strings to #hotdog.ColorRGBA
func MapCSSColor(colorString string) *hotdog.ColorRGBA {
	if string(colorString[0]) == "#" {
		return HexStringToColor(colorString)
	} else if rgba.MatchString(colorString) {
		return RGBAToColor(colorString)
	}

	return colorTable[colorString]
}
