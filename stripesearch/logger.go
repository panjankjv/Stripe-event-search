package stripesearch

import "log"

// DefaultLogger is default Logger.
var DefaultLogger Logger

// Logger is logging interface.
type Logger interface {
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

func init() {
	v := &DummyLogger{}
	DefaultLogger = v
}

// DummyLogger does not output anything
type DummyLogger struct{}

// Infof does nothing.
func (*DummyLogger) Infof(format string, v ...interface{}) {}

// Errorf does nothing.
func (*DummyLogger) Errorf(format string, v ...interface{}) {}

// StdLogger use standard log package.
type StdLogger struct{}

// Infof logging information.
func (*StdLogger) Infof(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

// Errorf logging error information.
func (*StdLogger) Errorf(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}
