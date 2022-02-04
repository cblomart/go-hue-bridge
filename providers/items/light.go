package items

var Lights map[int]*Light = make(map[int]*Light)

type Light struct {
	Provider string
	Name     string
	On       bool
	XID      string
}
