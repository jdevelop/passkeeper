package display

type DisplayControl interface {
	Refresh()
	ScrollUp(lines int)
	ScrollDown(lines int)
}
