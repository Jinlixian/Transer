package transer

//TransInput input
type TransInput struct {
	Query  string
	ID     string
	Secret string
	To     string
}

type FileTransIntput struct {
	Query      string
	ID         string
	Secret     string
	To         string
	InputFile  string
	OutputFile string
}

type TransOutput struct {
	Result string
}

//Transer interface for Trans
type Transer interface {
	Trans(input *TransInput) *TransOutput
}
