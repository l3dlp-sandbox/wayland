package mustard

import (
	"image"

	gg "github.com/danfragoso/thdwb/gg"


	"github.com/goki/freetype/truetype"
)

import window "github.com/neurlang/wayland/window"
import cairo "github.com/neurlang/wayland/cairoshim"

type Overlay struct {
	ref string

	active bool

	top  float64
	left float64

	width  float64
	height float64

	position image.Point

}

type contextMenu struct {
	overlay       *Overlay
	entries       []*menuEntry
	selectedEntry *menuEntry
}

type menuEntry struct {
	entryText string
	action    func()

	top  float64
	left float64

	width  float64
	height float64
}

type glBackend struct {
	program uint32

	vao uint32
	vbo uint32

	texture uint32
	quad    []float32
}

type box struct {
	top    float64
	left   float64
	width  float64
	height float64
}

func (box *box) SetCoords(top, left, width, height float64) {
	box.top = top
	box.left = left
	box.width = width
	box.height = height
}

func (box *box) GetCoords() (float64, float64, float64, float64) {
	return box.top, box.left, box.width, box.height
}

type Widget interface {
	SetNeedsRepaint(bool)
	NeedsRepaint() bool
	Widgets() []Widget
	ComputedBox() *box
	SetWindow(*Window)
	BaseWidget() *baseWidget

	render(s cairo.Surface, time uint32)
}

type baseWidget struct {
	box            box
	computedBox    box
	widgetPosition widgetPosition

	font *truetype.Font

	needsRepaint bool
	fixedWidth   bool
	fixedHeight  bool

	widgets []Widget

	backgroundColor string

	widgetType widgetType
	cursor     cursorType

	focusable  bool
	selectable bool

	focused  bool
	selected bool

	window *Window
	widget *window.Widget
}

type cursorType int

const (
	//DefaultCursor - Default arrow cursor
	DefaultCursor cursorType = 4
	//PointerCursor - Pointer cursor
	PointerCursor cursorType = 11
	//ArrowCursor - Arrow cursor
	ArrowCursor cursorType = 4
	//SpinnerCursor
	SpinnerCursor cursorType = 12
	//HandCursor
	HandCursor cursorType = 11
	//IbeamCursor
	IbeamCursor cursorType = 10

)

type widgetType int

const (
	buttonWidget widgetType = iota
	canvasWidget
	frameWidget
	imageWidget
	inputWidget
	labelWidget
	treeWidget
	scrollbarWidget
	textWidget
)

type widgetPosition int

const (
	PositionRelative widgetPosition = iota
	PositionAbsolute
)

type FrameOrientation int

const (
	//VerticalFrame - Vertical frame orientation
	VerticalFrame FrameOrientation = iota

	//HorizontalFrame - Horizontal frame orientation
	HorizontalFrame
)

type MustardKey int

const (
	MouseLeft MustardKey = iota
	MouseRight
)

//Frame - Layout frame type
type Frame struct {
	baseWidget

	orientation FrameOrientation
}

type LabelWidget struct {
	baseWidget
	content string

	fontSize  float64
	fontColor string
}

type TreeWidget struct {
	baseWidget

	fontSize  float64
	fontColor string
	nodes     []*TreeWidgetNode

	openIcon  image.Image
	closeIcon image.Image

	selectCallback func(*TreeWidgetNode)
}

func (widget *TreeWidget) RemoveNodes() {
	widget.nodes = nil
}

func (widget *TreeWidget) AddNode(childNode *TreeWidgetNode) {
	widget.nodes = append(widget.nodes, childNode)
}

func CreateTreeWidgetNode(key, value string) *TreeWidgetNode {
	return &TreeWidgetNode{
		Key:   key,
		Value: value,
		box:   box{},
	}
}

type TreeWidgetNode struct {
	Key      string
	Value    string
	Parent   *TreeWidgetNode
	Children []*TreeWidgetNode

	isOpen     bool
	isSelected bool
	index      int
	box        box
}

func (node *TreeWidgetNode) Toggle() {
	if node.isOpen {
		node.isOpen = false
	} else {
		node.isOpen = true
	}
}

func (node *TreeWidgetNode) Close() {
	node.isOpen = false

}
func (node *TreeWidgetNode) Open() {
	node.isOpen = true
}

func (node *TreeWidgetNode) AddNode(childNode *TreeWidgetNode) {
	childNode.Parent = node
	childNode.index = len(node.Children)
	node.Children = append(node.Children, childNode)
}

func (node *TreeWidgetNode) NextSibling() *TreeWidgetNode {
	selfIdx := node.index
	if selfIdx+1 < len(node.Parent.Children) {
		return node.Parent.Children[selfIdx+1]
	}

	return nil
}

func (node *TreeWidgetNode) PreviousSibling() *TreeWidgetNode {
	selfIdx := node.index
	if selfIdx-1 >= 0 {
		return node.Parent.Children[selfIdx-1]
	}

	return nil
}

type TextWidget struct {
	baseWidget
	content string

	fontSize  float64
	fontColor string
}

type ImageWidget struct {
	baseWidget

	path string
	img  image.Image
}

type CanvasWidget struct {
	baseWidget

	context        *gg.Context
	drawingContext *gg.Context

	renderer func(*CanvasWidget)

	scrollable bool
	offset     int

	drawingRepaint bool
}

type ButtonWidget struct {
	baseWidget
	content string

	icon      image.Image
	fontSize  float64
	fontColor string
	selected  bool
	padding   float64
	onClick   func()
}

type InputWidget struct {
	baseWidget

	value           string
	selected        bool
	active          bool
	padding         float64
	fontSize        float64
	context         *gg.Context
	fontColor       string
	cursorFloat     bool
	cursorPosition  int
	cursorDirection bool
	returnCallback  func()
}

type ScrollBarWidget struct {
	baseWidget

	orientation ScrollBarOrientation
	selected    bool
	thumbSize   float64
	thumbColor  string

	scrollerSize   float64
	scrollerOffset float64
}

type ScrollBarOrientation int

const (
	VerticalScrollBar ScrollBarOrientation = iota
	HorizontalScrollBar
)
