package entity

type User struct {
	ID          string `bson:"id"`
	Name        string `bson:"name"`
	Email       string `bson:"email"`
	Password    string `bson:"password"`
	IsConfirmed bool   `bson:"is_confirmed"`
}
