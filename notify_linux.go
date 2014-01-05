// Copyright (c) 2013, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

//+build linux bsd

package main

import (
	"fmt"

	"github.com/goulash/notify"
)

var notificationID uint32

func init() {
	if notify.ServiceAvailable() {
		notify.SetName("mywire")
		trayNotify = func(goal NotifyLevel, args ...interface{}) {
			msg := fmt.Sprint(args...)
			notificationID, _ = notify.ReplaceUrgentMsg(notificationID, msg, "", notify.NotificationUrgency(goal+1%2))
		}
	}
}
