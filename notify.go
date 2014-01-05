// Copyright (c) 2013, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package main

type NotifyLevel int

const (
	NotifyNone NotifyLevel = iota
	NotifyCritical
	NotifyNormal
	NotifyAll
)

func (nl NotifyLevel) TrayNotify(goal NotifyLevel, args ...interface{}) {
	if nl >= goal {
		trayNotify(goal, args...)
	}
}

var trayNotify = notifyNoOne

// notifyNoOne is a function that can be assigned to Notify that does nothing.
func notifyNoOne(NotifyLevel, ...interface{}) {}
