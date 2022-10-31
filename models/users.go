package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const UsersCollections = "users"

type UserModel struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	Phone          string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Email          string             `bson:"email,omitempty" json:"email,omitempty"`
	DOB            string             `bson:"dob" json:"dob"`
	OTP            string             `bson:"otp" json:"otp"`
	Password       string             `bson:"password" json:"password"`
	ProfilePicture string             `bson:"profile" json:"profile"`
	UserName       string             `bson:"handle" json:"handle"`
	IsVerified     bool               `bson:"is_verified" json:"is_verified"`
	RefreshToken   string             `bson:"refresh_token" json:"refreshToken"`
	Followers      int                `bson:"followers" json:"followers"`
	Following      int                `bson:"following" json:"following"`
}
