package logger

type ModeLogger string

const (
	ModeDev  ModeLogger = "dev"
	ModeProd ModeLogger = "prod"
)

func (m ModeLogger) IsDev() bool {
	return m == ModeDev
}

func (m ModeLogger) IsProd() bool {
	return m == ModeProd
}
