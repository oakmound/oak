package oak

// Need a way to mock os.Stdin to do this
//
// func TestDebugConsole(t *testing.T) {
// 	triggered := false
// 	AddCommand("test", func([]string) {
// 		triggered = true
// 	})
// 	rCh := make(chan bool)
// 	sCh := make(chan bool)
// 	//os.Stdin = bytes.NewBuffer([]byte("c test\n"))
// 	// stdinWriter := bufio.NewWriter(os.Stdin)
// 	// stdinWriter.WriteString("c test\n")
// 	// stdinWriter.Flush()
// 	go debugConsole(rCh, sCh)
// 	sleep()
// 	sleep()
// 	assert.True(t, triggered)
// }
