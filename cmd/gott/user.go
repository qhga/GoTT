package main

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// User of the experiment
type User struct {
	ID          primitive.ObjectID `bson:"_id"`
	UID         string             `bson:"uid"`
	Mail        string             `bson: mail`
	Pass        string             `bson:"pass"`
	CurrSession Session            `bson:"curr_session"`
	Role        string             `bson:"role"`
}

func init() {
	createUser(config.WebAdmin, "", config.WebAdminPass, "admin")
}

// createUser creates a new participant and writes it to the database.
func createUser(uid string, email string, pass string, role string) *User {
	_, err := getUserByUID(uid)
	if err != mongo.ErrNoDocuments {
		log.Println("User with that UID already exists")
		return nil
	}
	if len(email) > 0 {
		_, err = getUserByEmail(email)
		if err != mongo.ErrNoDocuments {
			log.Println("User with that email already exists")
			return nil
		}
	}
	// salt+hash password to store in DB
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	u := &User{
		ID:   primitive.NewObjectID(),
		UID:  uid,
		Mail: email,
		Pass: string(hashedPw),
		Role: role,
	}
	_, err = colUsers.InsertOne(mongoCtx, u)
	if err != nil {
		log.Fatal(err)
	}
	return u
}

// generateUserID generates a unique ID for each participant
func generateUserID() (uid string) {
	for {
		uid = generateRandomString(8, "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
		// Check if ID is already in use
		err := colUsers.FindOne(mongoCtx, bson.M{"uid": uid}).Decode(bson.M{})
		// If not in use stop loop
		if err == mongo.ErrNoDocuments {
			break
		} else {
			// Log other errors
			log.Fatal(err)
		}
	}
	return
}

// getUserByUID returns a user with the same UID as uid from the DB
func getUserByUID(uid string) (p *User, err error) {
	p = &User{}
	err = colUsers.FindOne(mongoCtx, bson.M{"uid": uid}).Decode(p)
	if err == mongo.ErrNoDocuments {
		return nil, err
	} else if err != nil {
		log.Fatal(err)
	}

	return p, nil
}

// getUserByEmail returns a user with the same email as email from the DB
func getUserByEmail(email string) (p *User, err error) {
	p = &User{}
	err = colUsers.FindOne(mongoCtx, bson.M{"mail": email}).Decode(p)
	if err == mongo.ErrNoDocuments {
		return nil, err
	} else if err != nil {
		log.Fatal(err)
	}
	return p, nil
}

// getUserBySession returns a user with the same session as session from the DB
func getUserBySession(s *Session) (p *User, err error) {
	p = &User{}

	// Get only id and authtoken otherwise timestamp might be already overwritten
	err = colUsers.FindOne(mongoCtx, bson.M{
		"curr_session.id":         s.ID,
		"curr_session.auth_token": s.AuthToken,
	}).Decode(p)
	if err == mongo.ErrNoDocuments {
		return nil, err
	} else if err != nil {
		log.Fatal(err)
	}
	return p, nil
}
