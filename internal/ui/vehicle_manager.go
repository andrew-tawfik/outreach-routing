package ui

import "fyne.io/fyne/v2"

type dragManager struct {
	overlay      *fyne.Container
	draggedGuest *GuestWidget
	originalPos  fyne.Position
	currentHover *GuestTile
	validTargets []*GuestTile
}
