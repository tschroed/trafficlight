package trafficlight

const (
	OFF = 0
	RED = 1
	GREEN = 2
	AMBER = 4
)

type TrafficLight interface {
	// It's expected that Set understands how to demux the word into the
	// necessary driver-specific commands.
	Set(word uint8) error
}
