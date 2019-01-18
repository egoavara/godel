package godel1

type Click uint16

const (
	CursorKeyPrimary   Click = iota
	CursorKeySecondary Click = iota
	CursorKeyMiddle    Click = iota
	CursorKeyPageUp    Click = iota
	CursorKeyPageDown  Click = iota
)
