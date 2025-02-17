package themes

import (
	"main/drawer_theme"
)

func CreateThemeDoom() *drawer_theme.Theme {
	return drawer_theme.CreateTheme(
		"#1c1f24", // background
		"#ff6c6b", // red
		"#98be65", // green
		"#da8548", // yellow
		"#c51605", // orange
		"#51afef", // blue
		"#c678dd", // magenta
		"#5699af", // cyan
		"#202328", // foreground
	)
}

func CreateThemeDracula() *drawer_theme.Theme {
	return drawer_theme.CreateTheme(
		"#21222c", // background
		"#ff5555", // red
		"#50fa7b", // green
		"#f1fa8c", // yellow
		"#e78a4e", // orange
		"#bd93f9", // purple
		"#ff79c6", // pink
		"#8be9fd", // cyan
		"#f8f8f2", // foreground
	)
}

func CreateThemeGruvBox() *drawer_theme.Theme {
	return drawer_theme.CreateTheme(
		"#1d2021", // background
		"#ea6962", // red
		"#a9b665", // green
		"#d8a657", // yellow
		"#e78a4e", // orange
		"#7daea3", // blue
		"#d3869b", // purple
		"#89b482", // cyan
		"#d4be98", // foreground
	)
}

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

func CreateThemeSweetMars() *drawer_theme.Theme {
	return drawer_theme.CreateTheme(
		"#1a1e21", // background
		"#cc241d", // red
		"#98971a", // green
		"#d79921", // yellow
		"#e78a4e", // orange
		"#458588", // blue
		"#b16286", // purple
		"#689d6a", // cyan
		"#a89984", // foreground
	)
}

func CreateThemeTomorrowDark() *drawer_theme.Theme {
	return drawer_theme.CreateTheme(
		"#1D1F21", // background
		"#a54242", // red
		"#8c9440", // green
		"#de935f", // yellow
		"#e78a4e", // orange
		"#5f819d", // blue
		"#85678f", // purple
		"#5e8d87", // cyan
		"#707880", // foreground
	)
}

func CreateThemeTomorrowNight() *drawer_theme.Theme {
	return drawer_theme.CreateTheme(
		"#1d1f21", // background
		"#cc6666", // red
		"#b5bd68", // green
		"#e6c547", // yellow
		"#e78a4e", // orange
		"#81a2be", // blue
		"#b294bb", // purple
		"#70c0ba", // cyan
		"#373b41", // foreground
	)
}

func CreateThemeTomorrow() *drawer_theme.Theme {
	return drawer_theme.CreateTheme(
		"#1D1F21", // background
		"#fb4934", // red
		"#98971a", // green
		"#d7992a", // yellow
		"#e78a4e", // orange
		"#458588", // blue
		"#b16286", // purple
		"#689d6a", // cyan
		"#373b41", // foreground
	)
}

// New Themes

func CreateThemeSolarizedDark() *drawer_theme.Theme {
	return drawer_theme.CreateTheme(
		"#002b36", // background (base03)
		"#dc322f", // red
		"#859900", // green
		"#b58900", // yellow
		"#cb4b16", // orange
		"#268bd2", // blue
		"#d33682", // magenta
		"#2aa198", // cyan
		"#839496", // foreground (base0)
	)
}

func CreateThemeSolarizedLight() *drawer_theme.Theme {
	return drawer_theme.CreateTheme(
		"#fdf6e3", // background (base3)
		"#dc322f", // red
		"#859900", // green
		"#b58900", // yellow
		"#cb4b16", // orange
		"#268bd2", // blue
		"#d33682", // magenta
		"#2aa198", // cyan
		"#93a1a1", // foreground (base1)
	)
}

func CreateThemeNord() *drawer_theme.Theme {
	return drawer_theme.CreateTheme(
		"#2E3440", // background (Nord0)
		"#BF616A", // red
		"#3B4252", // green
		"#434C5E", // yellow
		"#4C566A", // orange
		"#8FBCBB", // blue
		"#88C0D0", // magenta
		"#81A1C1", // cyan
		"#a9bbd1", // foreground (or accent)
	)
}

func CreateThemeMonokai() *drawer_theme.Theme {
	return drawer_theme.CreateTheme(
		"#272822", // background
		"#F92672", // red
		"#A6E22E", // green
		"#E6DB74", // yellow
		"#FD971F", // orange
		"#66D9EF", // blue
		"#AE81FF", // purple
		"#75715E", // cyan (used for comments/secondary)
		"#F8F8F2", // foreground
	)
}
