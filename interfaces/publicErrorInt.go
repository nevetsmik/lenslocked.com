package interfaces

type PublicErrorInt interface {
	error
	Public() string
}
