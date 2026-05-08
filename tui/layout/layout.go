package layout

const (
	InputHeight       = 5
	StatusBarHeight   = 1
	SpacingHeight     = 1
	PaddingHorizontal = 1
)

func GetChatWidth(totalWidth int) int {
	return totalWidth - PaddingHorizontal*2
}

func GetChatHeight(totalHeight int) int {
	return totalHeight - InputHeight - SpacingHeight - StatusBarHeight
}

func GetInputWidth(totalWidth int) int {
	return totalWidth - PaddingHorizontal*2
}

func GetInputHeight() int {
	return InputHeight
}

func GetStatusBarHeight() int {
	return StatusBarHeight
}

func GetSpacingHeight() int {
	return SpacingHeight
}

func GetPaddingHorizontal() int {
	return PaddingHorizontal
}