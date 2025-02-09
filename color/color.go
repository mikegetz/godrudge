package color

type Color string

const Blue Color = "\033[34m"
const Red Color = "\033[31m"

const End Color = "\033[0m"

func ColorString(color Color, message string) string {
	return string(color) + message + string(End)
}
