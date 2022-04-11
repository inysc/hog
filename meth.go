package qog

import "fmt"

// ************* TRACE ****************
func (l *Logger) Trace(msg string) {
	if l.lvl > TRACE {
		return
	}
	l.write(l.trace, msg)
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	if l.lvl > TRACE {
		return
	}
	msg := fmt.Sprintf(format, args...)
	l.write(l.trace, msg)
}

// ************* DEBUG ****************
func (l *Logger) Debug(msg string) {
	if l.lvl > DEBUG {
		return
	}
	l.write(l.debug, msg)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.lvl > DEBUG {
		return
	}
	msg := fmt.Sprintf(format, args...)
	l.write(l.debug, msg)
}

// ************* INFO ****************
func (l *Logger) Info(msg string) {
	if l.lvl > INFO {
		return
	}
	l.write(l.info, msg)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.lvl > INFO {
		return
	}
	msg := fmt.Sprintf(format, args...)
	l.write(l.info, msg)
}

// ************* Warn ****************
func (l *Logger) Warn(msg string) {
	if l.lvl > WARN {
		return
	}
	l.write(l.warn, msg)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.lvl > WARN {
		return
	}
	msg := fmt.Sprintf(format, args...)
	l.write(l.warn, msg)
}

// ************* ERROR ****************
func (l *Logger) Error(msg string) {
	if l.lvl > ERROR {
		return
	}
	l.write(l.err, msg)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.lvl > ERROR {
		return
	}
	msg := fmt.Sprintf(format, args...)
	l.write(l.err, msg)
}
