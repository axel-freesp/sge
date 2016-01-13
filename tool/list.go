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
	i, ok := l.Find(s)
	if !ok {
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

func (l StringList) Strings() []string {
	return l.strings
}

func (l StringList) Find(s string) (index int, ok bool) {
	for index = range l.strings {
		if s == l.strings[index] {
			ok = true
			break
		}
	}
	return
}
