package behaviour

import (
	"testing"
)

/*
type SignalType interface {
	TreeElement
	TypeName() string
	CType() string
	ChannelId() string
	Scope() Scope
	Mode() Mode
}
*/

func TestSignalType(t *testing.T) {
	case1 := []struct {
		remove, name, ctype, msgid string
		scope                      Scope
		mode                       Mode
		isLegal                    bool
	}{
		/*
		 *	If 'remove' is not empty, the signal shall be removed
		 * 	(all other fields are irrelevant).
		 * 	If 'remove' is empty, try to register the signal given by
		 * 	all the other definitions.
		 *
		 * 	'isLegal' is the expected behaviour: true: testcase shall
		 * 	be successful, false: testcase shall fail
		 *
		 * 	Compatible duplicates are allowed (all fields identical!),
		 * 	but duplicates are registered only once.
		 *
		 * 	TreeElement interface not tested.
		 */
		{"", "s1", "", "", 0, 0, true},
		{"", "s1", "", "", 0, 0, true},        // compatible duplicate
		{"", "s1", "int", "", 0, 0, false},    // incompatible duplicate
		{"", "s2", "int", "ch1", 0, 0, true},  // new
		{"", "s2", "int", "ch1", 0, 1, false}, // incompatible duplicate
		{"", "s2", "int", "ch1", 1, 0, false}, // incompatible duplicate
		{"", "s2", "int", "ch1", 0, 0, true},  // compatible duplicate
		{"s1", "", "", "", 0, 0, true},        // remove
		{"s1", "", "", "", 0, 0, false},       // remove duplicate
		{"s3", "", "", "", 0, 0, false},       // remove non-existing
		{"", "s1", "int", "", 0, 0, true},     // new (after removal)
	}
	Init()
	for i, c := range case1 {
		if len(c.remove) > 0 {
			s := signalTypes[c.remove]
			success := (s != nil)
			if success != c.isLegal {
				t.Errorf("TestSignalType testcase %d failed, could not find %s\n", c.remove)
			} else if success {
				SignalTypeDestroy(s)
			}
		} else {
			s, err := SignalTypeNew(c.name, c.ctype, c.msgid, c.scope, c.mode)
			success := (err == nil)
			if success != c.isLegal {
				t.Errorf("TestSignalType testcase %d failed, err=%v.\n", i, err)
			} else if success {
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
}
