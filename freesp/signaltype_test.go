package freesp

import (
	"testing"
)

func TestSignalType(t *testing.T) {
	case1 := []struct {
		remove, name, ctype, msgid string
		scope                      Scope
		mode                       Mode
		isLegal                    bool
	}{
		{"", "s1", "", "", 0, 0, true},
		{"", "s1", "", "", 0, 0, true},        // compatible duplicate
		{"", "s1", "int", "", 0, 0, false},    // incompatible duplicate
		{"", "s2", "int", "ch1", 0, 0, true},  // new
		{"", "s2", "int", "ch1", 0, 1, false}, // incompatible duplicate
		{"", "s2", "int", "ch1", 1, 0, false}, // incompatible duplicate
		{"", "s2", "int", "ch1", 0, 0, true},  // compatible duplicate
		{"s1", "s1", "int", "", 0, 0, true},   // new (after removal)
	}
	Init()
	for i, c := range case1 {
		if len(c.remove) > 0 {
			s := signalTypes[c.remove]
			if s == nil {
				t.Errorf("TestSignalType testcase %d failed, could not find %s\n", c.remove)
			} else {
				SignalTypeDestroy(s)
			}
		}
		s, err := SignalTypeNew(c.name, c.ctype, c.msgid, c.scope, c.mode)
		if (err == nil) != c.isLegal {
			t.Errorf("TestSignalType testcase %d failed, err=%v.\n", i, err)
		}
		if err == nil {
			if s.TypeName() != c.name {
				t.Errorf("TestSignalType testcase %d failed, TypeName().\n", i)
			}
			if s.CType() != c.ctype {
				t.Errorf("TestSignalType testcase %d failed, CType().\n", i)
			}
			if s.ChannelId() != c.msgid {
				t.Errorf("TestSignalType testcase %d failed, ChannelId().\n", i)
			}
			if s.Scope() != c.scope {
				t.Errorf("TestSignalType testcase %d failed, Scope().\n", i)
			}
			if s.Mode() != c.mode {
				t.Errorf("TestSignalType testcase %d failed, Mode().\n", i)
			}
		}
	}
}
