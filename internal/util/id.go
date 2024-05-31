package util

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenerateID() string {
	return gonanoid.Must()
}
