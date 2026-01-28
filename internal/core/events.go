package core

type Event struct {
	Type    string
	Payload any
}

const (
	EventMessage        = "message"
	EventPeerFound      = "peer_found"
	EventChatRequest    = "chat_request"
	EventChatAccept     = "chat_accept"
	EventChatReject     = "chat_reject"
	EventChatDisconnect = "chat_disconnect"
)
