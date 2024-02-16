package id

import gonanoid "github.com/matoous/go-nanoid/v2"

func Generate() string {
	id, err := gonanoid.New()
	if err != nil {
		panic("error generating id: " + err.Error())
	}

	return id
}
