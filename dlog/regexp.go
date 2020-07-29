package dlog

import (
	"regexp"
	"errors"
	"fmt"
)

type RegexpLogger interface {
	FullLogger
	SetRegexp(*regexp.Regexp)
}

type regexpLogger struct {
	*logger
	regexp *regexp.Regexp
	overrideFilterLevel Level
}

func NewRegexpLogger() RegexpLogger {
	return &regexpLogger{
		NewLogger().(*logger),
		nil,
		WARN,
	}
}

// addRegexp convertes a normal logger into a regexp logger.
// reg can be:
// - a string representing a regexp
// - an already-compiled regexp (or a pointer to one)
// - nil (no regexp)
func addRegexp(l *logger,reg interface{}) (regexpLogger, error) {
	var regex *regexp.Regexp
	if reg == nil {
		regex = nil
	}
	switch reg.(type) {
	case *regexp.Regexp:
		regex = reg.(*regexp.Regexp)
	case regexp.Regexp:
		r := reg.(regexp.Regexp)
		regex = &r
	case string:
		var err error
		regex, err = regexp.Compile(reg.(string))
		if err != nil {
			return regexpLogger{}, err
		}
	default:
		return regexpLogger{}, errors.New("invalid type")
	}
	return regexpLogger{l,regex,WARN}, nil
}

// SetOverrideLevel sets the log level needed to bypass the regex.
func (l *regexpLogger) SetOverrideLevel(lvl Level) {
	l.overrideFilterLevel = lvl
}

func (l *regexpLogger) SetRegexp(r *regexp.Regexp) {
	l.regexp = r
}

func (l *regexpLogger) SetDebugFilter(regexpStr string) {
	l.regexp = regexp.MustCompile(regexpStr)
}

func (l regexpLogger) Error(in ...interface{}) {
	l.logCond(in,ERROR)
}

func (l regexpLogger) Warn(in ...interface{}) {
	l.logCond(in,WARN)
}

func (l regexpLogger) Info(in ...interface{}) {
	l.logCond(in,INFO)
}

func (l regexpLogger) Verb(in ...interface{}) {
	l.logCond(in,VERBOSE)
}

func (l regexpLogger) logCond(in []interface{},lvl Level) {
	if l.debugLevel >= lvl &&
		(l.regexp == nil ||
			l.overrideFilterLevel >= lvl || l.regexp.MatchString(fmt.Sprint(in))) {

		l.logger.dLog(true, true, in)
	}
}
