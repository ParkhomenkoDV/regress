package export

const (
	before = "before"
	after  = "after"
)

type action string

const (
	change action = "change"
	add    action = "add"
	del    action = "del"
)
