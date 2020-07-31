package osm

type Logger interface {
	Infow(string, ...interface{})
	Debugw(string, ...interface{})
	Warnw(string, ...interface{})
	Error(...interface{})
	Errorw(string, ...interface{})
}
