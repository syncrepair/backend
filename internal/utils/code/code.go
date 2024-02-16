package code

import gonanoid "github.com/matoous/go-nanoid/v2"

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const length = 7

func Generate() string {
	id, err := gonanoid.Generate(alphabet, length)
	if err != nil {
		panic("error generating id: " + err.Error())
	}

	return id
}
