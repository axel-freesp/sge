package tool

import (
	"log"
)

type StringList struct {
	strings []string
}

func StringListInit() StringList {
	return StringList{nil}
}

func (l *StringList) Append(s string) {
	l.strings = append(l.strings, s)
}

func (l *StringList) Remove(s string) {
	var i int
	for i = range l.strings {
		if s == l.strings[i] {
			break
		}
	}
	if i >= len(l.strings) {
		for _, v := range l.strings {
			log.Printf("StringList.Remove: have %s\n", v)
		}
		log.Fatalf("StringList.Remove error: %s not in this list\n", s)
	}
	for i++; i < len(l.strings); i++ {
		l.strings[i-1] = l.strings[i]
	}
	l.strings = l.strings[:len(l.strings)-1]
}

func (l *StringList) Strings() []string {
	return l.strings
}
