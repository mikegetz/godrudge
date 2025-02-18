package color

type Color string

const Blue Color = "\033[34m"
const Red Color = "\033[31m"

const Reset Color = "\033[0m"

const hrefStart = "\033]8;;"
const hrefEnd = "\033\\"
const hrefTextEnd = "\033]8;;\033\\"

// adds Unicode color c to message string
func ColorString(c Color, s string) string {
	return string(c) + s + string(Reset)
}

func AnsiLink(href string, s string) string {
	return hrefStart + href + hrefEnd + s + hrefTextEnd
}
