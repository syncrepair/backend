package hasher

type Hasher interface {
	Hash(string) string
	Check(string, string) bool
}
