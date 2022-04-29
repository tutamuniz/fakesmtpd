package chat

type Chat interface {
	SendMessage(user, message string)
}
