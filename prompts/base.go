package prompts

type Template interface {
	Format(values map[string]string) string
	GetVariables() []string
}
