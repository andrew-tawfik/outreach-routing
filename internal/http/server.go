package http

func StartServer() string {
	return "server starting"
}

type SheetRepository interface {
	GetGuestAddresses()
}
