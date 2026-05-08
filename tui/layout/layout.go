package layout

const (
	RightWidth            = 40
	InputHeight           = 5
	StatusBarHeight       = 1
	SpacingHeight         = 1
	BarWidth              = 1
	LeftPaddingHorizontal = 1
	LeftPaddingVertical   = 1
)

type Dims struct {
	Width  int
	Height int
}

func GetLeftWidth(totalWidth int) int {
	return totalWidth - RightWidth - BarWidth
}

func GetRightWidth() int {
	return RightWidth
}

func GetChatHeight(totalHeight int) int {
	return totalHeight - InputHeight - SpacingHeight - StatusBarHeight
}

func GetInputWidth(totalWidth int) int {
	return totalWidth - LeftPaddingHorizontal*2
}

func GetChatWidth(totalWidth int) int {
	return totalWidth - BarWidth - LeftPaddingHorizontal*2
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

func GetBarWidth() int {
	return BarWidth
}

func GetLeftPaddingHorizontal() int {
	return LeftPaddingHorizontal
}

func GetLeftPaddingVertical() int {
	return LeftPaddingVertical
}
