package osm

type Logger interface {
	Info(...interface{})
	Infow(string, ...interface{})

	Debug(...interface{})
	Debugw(string, ...interface{})

	Warn(...interface{})
	Warnw(string, ...interface{})

	Error(...interface{})
	Errorw(string, ...interface{})
}
