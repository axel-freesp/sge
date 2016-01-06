package freesp

import (
	"testing"
)

func TestLibrary(t *testing.T) {
	case1 := []struct {
		library string
	}{
		{`<library xmlns="http://www.freesp.de/xml/freeSP" version="1.0">
   <signal-type name="s1" scope="" mode="" c-type="" message-id=""></signal-type>
   <node-type name="Test">
      <intype port="" type="s1"></intype>
      <outtype port="" type="s1"></outtype>
   </node-type>
</library>
`},
	}

	for i, c := range case1 {
		Init()
		var l Library = LibraryNew("test.alml", nil)
		buf := copyBuf(c.library)
		err := l.Read(buf)
		if err != nil {
			t.Errorf("Testcase %d: Failed to read from buffer: %v", i, err)
			return
		}
	}
}

func copyBuf(s string) (buf []byte) {
	buf = make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		buf[i] = s[i]
	}
	return
}
