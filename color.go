package godrudge

type Color string

const Blue Color = "\033[34m"
const Red Color = "\033[31m"

const End Color = "\033[0m"

// adds Unicode color c to message string
func colorString(c Color, message string) string {
	return string(c) + message + string(End)
}
