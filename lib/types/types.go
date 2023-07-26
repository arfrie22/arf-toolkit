package types

type ToolItem struct {
	Run         func()
	Name        string
	Description string
}
