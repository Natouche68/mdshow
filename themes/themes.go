package themes

import "github.com/charmbracelet/log"

type Theme struct {
	CodeStyle        string
	FontImportURL    string
	TextFontFamily   string
	TitleFontFamily  string
	CodeFontFamily   string
	BackgroundColor  string
	TextColor        string
	HUnderlined      bool
	H1Weight         int
	H1Color          string
	H2Weight         int
	H2Color          string
	H3Weight         int
	H3Color          string
	H4Weight         int
	H4Color          string
	H5Weight         int
	H5Color          string
	H6Weight         int
	H6Color          string
	DelColor         string
	StrongColor      string
	LinkColor        string
	LinkWeight       int
	TableBackground1 string
	TableBackground2 string
	TableTitleWeight int
	TableTitleColor  string
	CodeColor        string
	CodeBackground   string
	ArrowColor       string
	ProgressBarColor string
}

func GetTheme(themeName string) (Theme, string) {
	switch themeName {
	case "catppuccin":
		return CatppuccinTheme, "catppuccin"
	case "ocean":
		return OceanTheme, "ocean"
	default:
		log.Warn("Unknown theme, using default", "theme", themeName)
		return CatppuccinTheme, "catppuccin"
	}
}
