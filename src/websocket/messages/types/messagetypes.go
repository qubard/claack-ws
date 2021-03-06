package types

type MessageType = int8

const (
	InitSocket         MessageType = 0
	QueueRace          MessageType = 1
	UpdateRace         MessageType = 2
	UpdateRacer        MessageType = 3
	UpdateChat         MessageType = 4
	StatsUpdate        MessageType = 5
	ProfileUpdate      MessageType = 6
	ModalUpdate        MessageType = 7
	AddNotification    MessageType = 8
	RemoveNotification MessageType = 9
	UpdateWPM          MessageType = 10
	ClearRace          MessageType = 11
	UpdateRaceStart    MessageType = 12
	AddConversation    MessageType = 13
	RemoveConversation MessageType = 14
	AddMessage         MessageType = 15
	PlayerDisconnected MessageType = 16
	AuthUser           MessageType = 17 // Send in a token, authenticate user
	MessageDelivered   MessageType = 18
)

type RelayMessage struct {
	DstId   string // The id of the user we're sending the message to
	Message string // The packed message
}
