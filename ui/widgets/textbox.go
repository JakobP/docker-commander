package widgets

import (
	"image"

	. "github.com/gizak/termui/v3"
)

type TextBox struct {
	Block
	WrapText    bool
	TextStyle   Style
	CursorStyle Style
	ShowCursor  bool
	ID          string

	text        [][]Cell
	cursorPoint image.Point
}

var TextBoxTheme = TextBoxThemeType{
	Text:   NewStyle(ColorWhite),
	Cursor: NewStyle(ColorWhite, ColorClear, ModifierReverse),
}

type TextBoxThemeType struct {
	Text   Style
	Cursor Style
}

func NewTextBox() *TextBox {
	return &TextBox{
		Block:       *NewBlock(),
		WrapText:    false,
		TextStyle:   TextBoxTheme.Text,
		CursorStyle: TextBoxTheme.Cursor,

		text:        [][]Cell{[]Cell{}},
		cursorPoint: image.Pt(1, 1),
	}
}

func (self *TextBox) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	yCoordinate := 0
	for _, line := range self.text {
		if self.WrapText {
			line = WrapCells(line, uint(self.Inner.Dx()))
		}
		lines := SplitCells(line, '\n')
		for _, line := range lines {
			for _, cx := range BuildCellWithXArray(line) {
				x, cell := cx.X, cx.Cell
				buf.SetCell(cell, image.Pt(x, yCoordinate).Add(self.Inner.Min))
			}
			yCoordinate++
		}
		if yCoordinate > self.Inner.Max.Y {
			break
		}
	}

	if self.ShowCursor {
		point := self.cursorPoint.Add(self.Inner.Min).Sub(image.Pt(1, 1))
		cell := buf.GetCell(point)
		cell.Style = self.CursorStyle
		buf.SetCell(cell, point)
	}
}

func (self *TextBox) Backspace() {
	if self.cursorPoint == image.Pt(1, 1) {
		return
	}
	if self.cursorPoint.X == 1 {
		index := self.cursorPoint.Y - 1
		self.cursorPoint.X = len(self.text[index-1]) + 1
		self.text = append(
			self.text[:index-1],
			append(
				[][]Cell{append(self.text[index-1], self.text[index]...)},
				self.text[index+1:len(self.text)]...,
			)...,
		)
		self.cursorPoint.Y--
	} else {
		index := self.cursorPoint.Y - 1
		self.text[index] = append(
			self.text[index][:self.cursorPoint.X-2],
			self.text[index][self.cursorPoint.X-1:]...,
		)
		self.cursorPoint.X--
	}
}

// InsertText inserts the given text at the cursor position.
func (self *TextBox) InsertText(input string) {
	cells := ParseStyles(input, self.TextStyle)
	lines := SplitCells(cells, '\n')
	index := self.cursorPoint.Y - 1
	cellsAfterCursor := self.text[index][self.cursorPoint.X-1:]
	self.text[index] = append(self.text[index][:self.cursorPoint.X-1], lines[0]...)
	for i, line := range lines[1:] {
		index := self.cursorPoint.Y + i
		self.text = append(self.text[:index], append([][]Cell{line}, self.text[index:]...)...)
	}
	self.cursorPoint.Y += len(lines) - 1
	index = self.cursorPoint.Y - 1
	self.text[index] = append(self.text[index], cellsAfterCursor...)
	if len(lines) > 1 {
		self.cursorPoint.X = len(lines[len(lines)-1]) + 1
	} else {
		self.cursorPoint.X += len(lines[0])
	}
}

// ClearText clears the text and resets the cursor position.
func (self *TextBox) ClearText() {
	self.text = [][]Cell{[]Cell{}}
	self.cursorPoint = image.Pt(1, 1)
}

// SetText sets the text to the given text.
func (self *TextBox) SetText(input string) {
	self.ClearText()
	self.InsertText(input)
}

func (self *TextBox) GetText() string {
	var text string
	for _, r := range self.text {
		for _, t := range r {
			text += string(t.Rune)
		}
	}
	return text
}

func (self *TextBox) MoveCursorLeft() {
	self.MoveCursor(self.cursorPoint.X-1, self.cursorPoint.Y)
}

func (self *TextBox) MoveCursorRight() {
	self.MoveCursor(self.cursorPoint.X+1, self.cursorPoint.Y)
}

func (self *TextBox) MoveCursorUp() {
	self.MoveCursor(self.cursorPoint.X, self.cursorPoint.Y-1)
}

func (self *TextBox) MoveCursorDown() {
	self.MoveCursor(self.cursorPoint.X, self.cursorPoint.Y+1)
}

func (self *TextBox) MoveCursor(x, y int) {
	self.cursorPoint.Y = MinInt(MaxInt(1, y), len(self.text))
	self.cursorPoint.X = MinInt(MaxInt(1, x), len(self.text[self.cursorPoint.Y-1])+1)
}
