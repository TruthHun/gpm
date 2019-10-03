package conf

type GPM struct {
	WatchExt  []string
	WatchPath []string
	Commands  [][]string
	Frequency int
	Strict    bool
}

var Config GPM
