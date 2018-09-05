package protocol

const (
	HomeConsoleVersion = "Homeconsole02.00"
)

type Message interface {
	ParseContent(input string, pos int) TLV
}


