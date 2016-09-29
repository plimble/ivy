package ivy

type Source interface {
	Get(path string) ([]byte, error)
}
