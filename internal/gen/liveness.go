// Copyright (c) 2023 Justus Zorn

package gen

type liveness struct {
	parent *liveness
	vars   map[string]bool
}

func newLiveness(parent *liveness) *liveness {
	return &liveness{parent: parent, vars: map[string]bool{}}
}

func (l *liveness) isAlive(name string) bool {
	alive, ok := l.vars[name]
	if ok {
		// variable belongs to this liveness
		return alive
	} else if l.parent != nil {
		// ask parent liveness
		return l.parent.isAlive(name)
	} else {
		return false
	}
}

func (l *liveness) isOwned(name string) bool {
	// check whether variable belongs to this liveness scope
	_, ok := l.vars[name]
	return ok
}

func (l *liveness) makeAlive(name string) {
	l.vars[name] = true
}

func (l *liveness) makeDead(name string) {
	l.vars[name] = false
}

func (l *liveness) cloneLiveness() *liveness {
	new := newLiveness(l.parent)
	for name, alive := range l.vars {
		new.vars[name] = alive
	}
	return new
}

func mergeLiveness(s *scope, left *liveness, right *liveness) *liveness {
	new := newLiveness(left.parent)
	for name := range s.vars {
		if left.isAlive(name) && right.isAlive(name) {
			new.vars[name] = true
		} else {
			new.vars[name] = false
		}
	}
	return new
}
