package models

type User struct {
	ID       int64
	Email    string
	PassHash []byte
}

const EmptyUserID = int64(0)
