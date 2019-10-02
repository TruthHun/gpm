package conf

type GPM struct {
	WatchExt  []string
	WatchPath []string
	Commands  []string
	Frequency int
}

var Config GPM
