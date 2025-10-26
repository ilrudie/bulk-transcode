package ffmpeg

type Args struct {
	Options []string `yaml:"options"`
}

func (a *Args) Generate(i, o string) string {
	return "not implemented"
}
