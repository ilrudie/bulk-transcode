package ffmpeg

type Exec struct {
	Input, Output string
}

func New(i, o string) *Exec {
	return &Exec{
		Input:  i,
		Output: o,
	}
}

func (e *Exec) Generate(args Args) string {
	cmd := "ffmpeg "
	for _, option := range args.InputOptions {
		cmd += option + " "
	}
	cmd += "-i " + e.Input + " "
	for _, option := range args.OutputOptions {
		cmd += option + " "
	}
	cmd += e.Output
	return cmd
}
