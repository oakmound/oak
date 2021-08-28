// Package dlog provides basic logging functions
//
// It is not intended to be a fully featured or fully optimized logger-- it is
// just enough of a logger for oak's needs. A program utilizing oak, if it wants
// more powerful logs, should log to a more powerful tool, and if desired, tell oak
// to as well via setting dlog.DefaultLogger.
package dlog
