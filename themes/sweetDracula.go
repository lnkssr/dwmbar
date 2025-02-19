package themes

import (
	"main/drawer_theme"
)

func CreateThemeSweetDracula() *drawer_theme.Theme {
	return drawer_theme.CreateTheme(
		"#1a1e21", // background
		"#ff6c6b", // red
		"#50fa7b", // green
		"#f1fa8c", // yellow
		"#e78a4e", // orange
		"#bd93f9", // purple
		"#ff85ab", // pink (changed to avoid duplicate red)
		"#8be9fd", // cyan
		"#f8f8f2", // foreground
	)
}
