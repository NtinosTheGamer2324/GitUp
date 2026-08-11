package helper

const (
	colorReset   = "\033[0m"
	colorBold    = "\033[1m"
	colorDim     = "\033[2m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[97m"
)

func wrap(code, s string) string {
	return code + s + colorReset
}

func Red(s string) string     { return wrap(colorRed, s) }
func Green(s string) string   { return wrap(colorGreen, s) }
func Yellow(s string) string  { return wrap(colorYellow, s) }
func Blue(s string) string    { return wrap(colorBlue, s) }
func Magenta(s string) string { return wrap(colorMagenta, s) }
func Cyan(s string) string    { return wrap(colorCyan, s) }
func White(s string) string   { return wrap(colorWhite, s) }
func Bold(s string) string    { return wrap(colorBold, s) }
func Dim(s string) string     { return wrap(colorDim, s) }
