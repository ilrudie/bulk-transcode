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
	for _, option := range args.Options {
		cmd += option + " "
	}
	cmd += "-i " + e.Input + " " + e.Output
	return cmd
}