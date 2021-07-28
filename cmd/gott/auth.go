package main

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	ID        [32]byte           `json:"id" bson:"id"`
	AuthToken [32]byte           `json:"auth_token" bson:"auth_token"`
	Timestamp int64              `json:"timestamp" bson:"timestamp"`
	LastTTID  primitive.ObjectID `json:"last_ttid" bson:"last_ttid"`
	LastKID   string             `json:"last_kid" bson:"last_kid"`
}

// ------------------------------- MIDDLEWARE ----------------------------------
type ctxKey int

const currUserKey ctxKey = 0

// validateSession verifies that the client has a session that is still valid
func validateSession(h http.HandlerFunc, role string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Allow to access required files pre login
		if r.URL.Path == "/static/css/tt.css" ||
			r.URL.Path == "/static/js/helper.js" ||
			r.URL.Path == "/static/fonts/OpenSans-Regular.ttf" {
			h(w, r)
			return
		}

		cookie, err := r.Cookie("TTSESSION")

		if err == http.ErrNoCookie {
			// Redirect to /login
			w.Header().Set("Location", "/login?r="+r.URL.Path)
			w.WriteHeader(http.StatusFound)
			log.Println(err)
			return
		} else if err != nil {
			log.Fatal(err)
			return
		}

		// TESTING: Show cookie content
		// log.Printf("%+v", *cookie)
		requestedUser := checkSession([]byte(cookie.Value))

		if requestedUser == nil {
			unsetCookie(w)
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusFound)
			return
		}
		// Set new timestamp
		requestedUser.CurrSession.Timestamp = time.Now().UTC().UnixNano()
		// Save to database
		updateSession(requestedUser)

		setEncryptedCookie(w, requestedUser)
		ctx := context.WithValue(r.Context(), currUserKey, *requestedUser)
		requestWithUser := r.WithContext(ctx)

		// ORIGINAL HANDLER
		switch role {
		case "":
			h(w, requestWithUser)
		case "admin":
			if isAdmin(requestedUser) {
				h(w, requestWithUser)
			} else {
				w.Header().Set("Location", "/")
				w.WriteHeader(http.StatusFound)
			}
		}
	})
}

func isAdmin(u *User) bool {
	if u.Role == "admin" {
		log.Println("IS ADMIN")
		return true
	} else {
		log.Println("IS NOT ADMIN")
		return false
	}
}

func encryptSession(s Session) (es []byte) {
	js, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	es, err = encrypt(js, encKey)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func setEncryptedCookie(w http.ResponseWriter, u *User) {
	cValue := encryptSession(u.CurrSession)
	cookie := &http.Cookie{
		Name:     "TTSESSION",
		Value:    string(cValue),
		MaxAge:   int(6 * time.Hour),
		Path:     "/",
		HttpOnly: false,
	}
	http.SetCookie(w, cookie)
}

func unsetCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "TTSESSION",
		MaxAge: -69,
	}
	http.SetCookie(w, cookie)
}

// checkCredentials compares the pw provided by the user to the hashed version
// in the database and returns true for success and false otherwise
func checkCredentials(u *User, pw string) bool {
	hashedPw := u.Pass
	err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(pw))
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

// checkSession compares the encrypted session string with the current Session
// set for the user and returns true if they match
func checkSession(receivedSession []byte) *User {
	// Test
	dec, err := decrypt(receivedSession, encKey)
	if err != nil {
		log.Println("Encryption key changed, deleting remaining cookies!")
		return nil
	}
	// log.Println(string(dec))
	decryptedSession := Session{}
	err = json.Unmarshal(dec, &decryptedSession)
	if err != nil {
		log.Fatal(err)
	}

	// There is no current user but a (mayb) valid session -> get user with session from db
	requestedUser, err := getUserBySession(&decryptedSession)
	if err != nil {
		log.Println("No user with that session found in db:", err)
		return nil
	}

	// 6 hours ago
	maxSessionAge := time.Now().UTC().UnixNano() - int64(6*time.Hour)

	// Check if session is still valid (Other values must match anyways)
	if requestedUser.CurrSession.Timestamp >= maxSessionAge {
		return requestedUser
	} else {
		log.Println("User logged out because Session is no longer valid")
		return nil
	}
}

func updateSession(u *User) {
	filter := bson.D{{"uid", u.UID}}
	update := bson.D{{"$set", bson.D{{"curr_session", u.CurrSession}}}}
	res, err := colUsers.UpdateOne(mongoCtx, filter, update)
	if err != nil {
		log.Fatal(err, res)
	}
}

func createSession() (s Session) {
	var sID, aToken [32]byte
	_, err := io.ReadFull(cryptorand.Reader, sID[:])
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.ReadFull(cryptorand.Reader, aToken[:])
	if err != nil {
		log.Fatal(err)
	}
	// MAYB: Also hash the random number? Any advantages?
	// sID := sha256.Sum256(nonce)
	// aToken := sha256.Sum256(nonce)
	return Session{ID: sID, AuthToken: aToken, Timestamp: time.Now().UTC().UnixNano()}
}
