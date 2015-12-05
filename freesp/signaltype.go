package freesp

// signalType

type signalType struct {
	name, ctype, msgid string
	scope              Scope
	mode               Mode
}

func newSignalType(name, ctype, msgid string, scope Scope, mode Mode) *signalType {
	return &signalType{name, ctype, msgid, scope, mode}
}

func (t *signalType) TypeName() string {
	return t.name
}

func (t *signalType) CType() string {
	return t.ctype
}

func (t *signalType) ChannelId() string {
	return t.msgid
}

func (t *signalType) Scope() Scope {
	return t.scope
}

func (t *signalType) Mode() Mode {
	return t.mode
}
