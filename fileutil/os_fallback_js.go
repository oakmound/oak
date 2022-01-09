
package fileutil 

//go:build js
// +build js

func init() {
	// OS calls always fall in JS, disable calling to it by default 
	OSFallback = false 
}