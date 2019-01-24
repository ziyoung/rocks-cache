package cache

// Cache we need to achieve
type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
}
