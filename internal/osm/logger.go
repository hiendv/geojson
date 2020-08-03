package osm

type Logger interface {
	Infow(string, ...interface{})
	Debugw(string, ...interface{})
	Warnw(string, ...interface{})
	Error(...interface{})
	Errorw(string, ...interface{})
}

type LoggerEmpty struct{}

func (logger *LoggerEmpty) Infow(string, ...interface{})  {}
func (logger *LoggerEmpty) Debugw(string, ...interface{}) {}
func (logger *LoggerEmpty) Warnw(string, ...interface{})  {}
func (logger *LoggerEmpty) Error(...interface{})          {}
func (logger *LoggerEmpty) Errorw(string, ...interface{}) {}

var LoggerNoop = &LoggerEmpty{}
