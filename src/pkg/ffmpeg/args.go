package ffmpeg

type Args struct {
	InputOptions  []string `json:"input_options"`
	OutputOptions []string `json:"output_options"`
}

func (a *Args) Generate(i, o string) string {
	return "not implemented"
}
