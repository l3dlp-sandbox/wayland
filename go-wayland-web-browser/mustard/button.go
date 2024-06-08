package mustard

import (
	assets "github.com/neurlang/wayland/go-wayland-web-browser/assets"
	gg "github.com/danfragoso/thdwb/gg"
	"github.com/goki/freetype/truetype"
)

// CreateButtonWidget - Creates and returns a new Button Widget
func CreateButtonWidget(label string, asset []byte) *ButtonWidget {
	var widgets []Widget
	font, _ := truetype.Parse(assets.OpenSans(400))

	icon, _ := gg.LoadAsset(asset)

	return &ButtonWidget{
		baseWidget: baseWidget{
			needsRepaint: true,
			widgets:      widgets,

			widgetType:      buttonWidget,
			cursor:          HandCursor,
			backgroundColor: "#fff",

			font: font,
		},

		icon:      icon,
		content:   label,
		fontSize:  20,
		padding:   0,
		fontColor: "#000",
		selected:  false,
	}

}

// SetWidth - Sets the button width
func (button *ButtonWidget) SetWidth(width float64) {
	button.box.width = width
	button.fixedWidth = true
	button.RequestReflow()
}

// SetHeight - Sets the button height
func (button *ButtonWidget) SetHeight(height float64) {
	button.box.height = height
	button.fixedHeight = true
	button.RequestReflow()
}

// SetFontSize - Sets the button font size
func (button *ButtonWidget) SetFontSize(fontSize float64) {
	button.fontSize = fontSize
	button.needsRepaint = true
}

// SetFontSize - Sets the button font size
func (button *ButtonWidget) SetPadding(padding float64) {
	button.padding = padding
	button.needsRepaint = true
}

// SetContent - Sets the button content
func (button *ButtonWidget) SetContent(content string) {
	button.content = content
	button.needsRepaint = true
}

func (button *ButtonWidget) Click() {
	button.onClick()
}

// GetContent - Gets the button content
func (button *ButtonWidget) GetContent() string {
	return button.content
}

// SetFontColor - Sets the button font color
func (button *ButtonWidget) SetFontColor(fontColor string) {
	if len(fontColor) > 0 && string(fontColor[0]) == "#" {
		button.fontColor = fontColor
		button.needsRepaint = true
	}
}

// SetBackgroundColor - Sets the button background color
func (button *ButtonWidget) SetBackgroundColor(backgroundColor string) {
	if len(backgroundColor) > 0 && string(backgroundColor[0]) == "#" {
		button.backgroundColor = backgroundColor
		button.needsRepaint = true
	}
}

func (button *ButtonWidget) render(s Surface, time uint32) {
	context := makeContextFromCairo(s)
	top, left, width, height := button.computedBox.GetCoords()
	_, _ = width, height

	if button.selected {
		context.SetHexColor("#ccc")
	} else {
		context.SetHexColor("#ddd")
	}

	context.DrawRectangle(
		float64(left)+button.padding,
		float64(top)+button.padding,
		float64(width)-(button.padding*2),
		float64(height)-(button.padding*2),
	)
	context.Fill()

	if button.content != "" {
		context.SetHexColor(button.fontColor)
		context.SetFont(button.font, button.fontSize)
		context.DrawString(button.content, float64(left)+button.padding, float64(top)+button.padding+button.fontSize)
	}

	if button.icon != nil {
		context.DrawImage(button.icon, int(left+4), int(top+2))
	}
}
