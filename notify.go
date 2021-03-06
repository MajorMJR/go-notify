// Package notify provides an implementation of the Freedesktop Notifications Specification
// using the DBus API.
package notify

import "github.com/godbus/dbus"

// Notification expire timeout
const (
	ExpiresDefault = -1
	ExpiresNever   = 0
)

// Notification Categories
const (
	ClassDevice              = "device"
	ClassDeviceAdded         = "device.added"
	ClassDeviceError         = "device.error"
	ClassDeviceRemoved       = "device.removed"
	ClassEmail               = "email"
	ClassEmailArrived        = "email.arrived"
	ClassEmailBounced        = "email.bounced"
	ClassIm                  = "im"
	ClassImError             = "im.error"
	ClassImReceived          = "im.received"
	ClassNetwork             = "network"
	ClassNetworkConnected    = "network.connected"
	ClassNetworkDisconnected = "network.disconnected"
	ClassNetworkError        = "network.error"
	ClassPresence            = "presence"
	ClassPresenceOffline     = "presence.offline"
	ClassPresenceOnline      = "presence.online"
	ClassTransfer            = "transfer"
	ClassTransferComplete    = "transfer.complete"
	ClassTransferError       = "transfer.error"
)

// Urgency Levels
const (
	UrgencyLow      = byte(0)
	UrgencyNormal   = byte(1)
	UrgencyCritical = byte(2)
)

// Hints
const (
	HintActionIcons   = "action-icons"
	HintCategory      = "category"
	HintDesktopEntry  = "desktop-entry"
	HintImageData     = "image-data"
	HintImagePath     = "image-path"
	HintResident      = "resident"
	HintSoundFile     = "sound-file"
	HintSoundName     = "sound-name"
	HintSuppressSound = "suppress-sound"
	HintTransient     = "transient"
	HintX             = "x"
	HintY             = "y"
	HintUrgency       = "urgency"
)

type Capabilities struct {
	// Supports using icons instead of text for displaying actions.
	ActionIcons bool

	// The server will provide any specified actions to the user.
	Actions bool

	// Supports body text. Some implementations may only show the summary.
	Body bool

	// The server supports hyperlinks in the notifications.
	BodyHyperlinks bool

	// The server supports images in the notifications.
	BodyImages bool

	// Supports markup in the body text.
	BodyMarkup bool

	// The server will render an animation of all the frames in a given image array.
	IconMulti bool

	// Supports display of exactly 1 frame of any given image array.
	IconStatic bool

	// The server supports persistence of notifications. Notifications will be
	// retained until they are acknowledged or removed by the user or recalled by the sender.
	Persistence bool

	// The server supports sounds on notifications.
	Sound bool
}

// GetCapabilities returns the capabilities of the notification server.
func GetCapabilities() (c Capabilities, err error) {
	connection, err := dbus.SessionBus()
	if err != nil {
		return
	}
	obj := connection.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	call := obj.Call("org.freedesktop.Notifications.GetCapabilities", 0)
	if call.Err != nil {
		return Capabilities{}, call.Err
	}

	s := []string{}
	if err = call.Store(&s); err != nil {
		return
	}
	for _, v := range s {
		switch v {
		case "action-icons":
			c.ActionIcons = true
			break
		case "actions":
			c.Actions = true
			break
		case "body":
			c.Body = true
			break
		case "body-hyperlinks":
			c.BodyHyperlinks = true
			break
		case "body-images":
			c.BodyImages = true
			break
		case "body-markup":
			c.BodyMarkup = true
			break
		case "icon-multi":
			c.IconMulti = true
			break
		case "icon-static":
			c.IconStatic = true
			break
		case "persistence":
			c.Persistence = true
			break
		case "sound":
			c.Sound = true
			break
		}
	}
	return
}

type ServerInformation struct {
	// The name of the notification server daemon
	Name string

	// The vendor of the notification server
	Vendor string

	// Version of the notification server
	Version string

	// Spec version the notification server conforms to
	SpecVersion string
}

// GetServerInformation returns information about the notification server such as its name
// and version.
func GetServerInformation() (i ServerInformation, err error) {
	connection, err := dbus.SessionBus()
	if err != nil {
		return
	}
	obj := connection.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	call := obj.Call("org.freedesktop.Notifications.GetServerInformation", 0)
	if call.Err != nil {
		return ServerInformation{}, call.Err
	}

	if err = call.Store(&i.Name, &i.Vendor, &i.Version, &i.SpecVersion); err != nil {
		return
	}
	return
}

type Notification struct {
	// The optional name of the application sending the notification. Can be blank.
	AppName string

	// The optional notification ID that this notification replaces.
	ReplacesID uint32

	// The optional program icon of the calling application.
	AppIcon string

	// The summary text briefly describing the notification.
	Summary string

	// The optional detailed body text.
	Body string

	// The actions send a request message back to the notification client when invoked.
	Actions []string

	// Hints are a way to provide extra data to a notification server.
	Hints map[string]interface{}

	// The timeout time in milliseconds since the display of the notification at which
	// the notification should automatically close.
	Timeout int32
}

// NewNotification creates a new notification object with some basic information.
func NewNotification(summary, body string) Notification {
	return Notification{Summary: summary, Body: body, Timeout: ExpiresDefault}
}

// Show sends the information in the notification object to the server to be displayed.
func (n Notification) Show() (id uint32, err error) {

	// We need to convert the interface type of the map to dbus.Variant as people
	// dont want to have to import the dbus package just to make use of the notification
	// hints.
	hints := map[string]dbus.Variant{}
	if len(n.Hints) != 0 {
		for k, v := range n.Hints {
			hints[k] = dbus.MakeVariant(v)
		}
	}

	connection, err := dbus.SessionBus()
	if err != nil {
		return
	}
	obj := connection.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	call := obj.Call(
		"org.freedesktop.Notifications.Notify",
		0,
		n.AppName,
		n.ReplacesID,
		n.AppIcon,
		n.Summary,
		n.Body,
		n.Actions,
		hints,
		n.Timeout)
	if call.Err != nil {
		return 0, call.Err
	}

	if err = call.Store(&id); err != nil {
		return
	}
	return
}

// CloseNotification closes the notification if it exists using its id.
func CloseNotification(id uint32) (err error) {
	connection, err := dbus.SessionBus()
	if err != nil {
		return
	}
	obj := connection.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	call := obj.Call("org.freedesktop.Notifications.CloseNotification", 0, id)
	if call.Err != nil {
		return call.Err
	}
	return
}
