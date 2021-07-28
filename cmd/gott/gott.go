// author: phga
// backend for the portfolio
// docker-compose up -d
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"path"

	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// cheap routing
	http.HandleFunc("/", validateSession(handleIndex, ""))
	http.HandleFunc("/test", validateSession(handleTest, ""))
	http.HandleFunc("/text", validateSession(handleText, ""))
	http.HandleFunc("/fre", validateSession(handleFRE, ""))
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", validateSession(handleLogout, ""))
	http.HandleFunc("/survey", validateSession(handleSurvey, ""))
	http.HandleFunc("/result", validateSession(handleResult, ""))
	http.HandleFunc("/question", validateSession(handleQuestion, "admin"))
	http.HandleFunc("/user", validateSession(handleUser, "admin"))
	http.HandleFunc("/admin_panel", validateSession(handleAdminPanel, "admin"))

	// provide the inc directory to the useragent
	http.HandleFunc("/static/", validateSession(handleStatic, ""))
	// listen on port 8080 (I use nginx to proxy this local server)
	log.Fatalln(http.ListenAndServe(":"+config.WebPort, nil))
}

// ---------------------------- HELPER TYPES -----------------------------------

type TmplData struct {
	HeaderData interface{}
	BodyData   interface{}
}

// ------------------------- REQUEST HANDLING ----------------------------------

func handleStatic(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(http.Dir("../../web/static"))).ServeHTTP(w, r)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// Parses all required html files to provide the actual html that is shipped.
	t, _ := getTemplate("layouts/base.html", "tt.html", "header.html")
	currUser := r.Context().Value(currUserKey).(User)
	td := TmplData{currUser, struct {
		KB []Keyboard
		TT []string
	}{
		config.AvailableKeyboards,
		typingTestTextTTIDs,
	}}
	t.Execute(w, td)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := getTemplate("layouts/base.html", "login.html")

		lastURL := r.URL.Query().Get("r")
		if len(lastURL) == 0 {
			lastURL = "/"
		}
		td := TmplData{nil, struct{ RedirectTo string }{lastURL}}

		err := t.Execute(w, td)
		if err != nil {
			log.Println(err)
		}

	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Fatal(err)
		}
		uid := r.PostForm.Get("username")
		pw := r.PostForm.Get("password")

		// TESTING: Logging of entered credentials
		log.Println("UID: ", uid, " PW: ", pw)

		var requestedUser *User
		var err error
		requestedUser, err = getUserByUID(uid)
		if err != nil {
			log.Println("No user with this id exists: ", uid)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		success := checkCredentials(requestedUser, pw)
		if success {
			// IT IS IMPORTANT TO SET THE COOKIE BEFORE ANY OTHER OUTPUT TO W OCCURS
			requestedUser.CurrSession = createSession()
			updateSession(requestedUser)
			setEncryptedCookie(w, requestedUser)
			w.WriteHeader(http.StatusOK)
			return
		} else {
			log.Println("Failed login attempt for user: ", uid)
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	unsetCookie(w)
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}

func handleSurvey(w http.ResponseWriter, r *http.Request) {
	currUser := r.Context().Value(currUserKey).(User)
	lastKID := currUser.CurrSession.LastKID
	lastTTID := currUser.CurrSession.LastTTID
	if r.Method == http.MethodGet {
		requestedSurvey := r.URL.Query().Get("s")
		surveyQuestions, exists := config.Surveys[requestedSurvey]
		if !exists {
			log.Println("Survey with that identifier does not exist: ", requestedSurvey)
			return
		}
		survey := Survey{
			ID:   primitive.NewObjectID(),
			SID:  requestedSurvey,
			UID:  currUser.UID,
			KID:  lastKID,
			TTID: lastTTID,
		}

		survey.setQuestions(surveyQuestions...)

		t, _ := getTemplate("layouts/base.html", "survey.html", "header.html")
		td := TmplData{currUser, survey}
		t.Execute(w, td)
	} else if r.Method == http.MethodPost {
		survey := &Survey{}
		err := json.NewDecoder(r.Body).Decode(survey)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		switch survey.SID {
		case "demographics survey":
			err = saveSurvey(survey)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case "post typing test survey":
			if len(survey.KID) == 0 || survey.KID != lastKID {
				log.Println("Survey without or with wrong keyboard submitted! Aborting!")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Survey without or with wrong keyboard submitted!"}`))
				return
			}
			if len(survey.TTID) == 0 || survey.TTID != lastTTID {
				log.Println("Survey without or with wrong Test ID submitted! Aborting!")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "Survey without or with wrong Test ID submitted!"}`))
				return
			}
			if survey.UID != currUser.UID {
				log.Println("User tried to submit Survey for another user")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "User tried to submit Survey for another user"}`))
				return
			}
			for _, q := range survey.Questions {
				if len(q.Answer) == 0 {
					log.Println("Not all Questions have been answerd.")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error": "Not all Questions have been answerd!"}`))
					return
				}
			}
			err = saveSurvey(survey)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			err = saveSurvey(survey)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}

func handleResult(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := getTemplate("layouts/base.html", "result.html", "header.html")

		currUser := r.Context().Value(currUserKey).(User)
		var tts *[]TypingTest
		var err error
		if isAdmin(&currUser) {
			tts, err = getTypingTestResults(nil)
		} else {
			tts, err = getTypingTestResults(&currUser)
		}
		if err != nil {
			log.Fatal("Could not show any test results: ", err)
		}
		if len(*tts) <= 0 {
			tts = nil
		}
		td := TmplData{currUser, tts}
		t.Execute(w, td)
	}
}

func handleQuestion(w http.ResponseWriter, r *http.Request) {
	qDir := "../../web/data/questions"
	if r.Method == http.MethodGet {
		t, _ := getTemplate("layouts/base.html", "question.html", "header.html")
		currUser := r.Context().Value(currUserKey).(User)
		td := TmplData{currUser, struct{ Qtypes []string }{[]string{QtypeSingle, QtypeMulti, QtypeFreeText, QtypeAnalogSlider}}}
		t.Execute(w, td)
	} else if r.Method == http.MethodPost {
		q := &Question{QID: primitive.NewObjectID()}
		err := json.NewDecoder(r.Body).Decode(q)
		if err != nil {
			log.Println(err)
		}
		addNewQuestion(q, qDir)
	} else if r.Method == http.MethodPatch {
		scanDirForNewQuestions(qDir)
	}
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	currUser := r.Context().Value(currUserKey).(User)
	if r.Method == http.MethodGet {
		t, _ := getTemplate("layouts/base.html", "user.html", "header.html")
		td := TmplData{currUser, nil}
		t.Execute(w, td)
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Println(err, "Could not create User - Parsing of form failed.")
		}
		// Create new User and show some Values
		uid := r.PostForm.Get("username")
		email := r.PostForm.Get("email")
		pw := r.PostForm.Get("password")
		// Save clear text password for later use
		if len(uid) == 0 {
			uid = generateUserID()
		}
		if len(pw) < 8 {
			pw = generatePassword()
		}
		// Create a new user
		newUser := createUser(uid, email, pw, "")
		if newUser == nil {
			w.Header().Set("Location", "/user")
			// http.StatusFound redirects POST to GET otherwise this ends up in a loop
			w.WriteHeader(http.StatusFound)
			return
		}

		if len(email) != 0 {
			// send mail
			// TODO: Place email text into html template?
			subject := "Your Account is ready"
			mailTxt := `Hey,

thanks for participating in my study about efficiency and satisfaction on different keyboards.

Here are your login credentials:

Username: %s
Password: %s

Cheers and see you soon,

Philip`
			mailTxt = fmt.Sprintf(mailTxt, uid, pw)
			sendMail(email, subject, mailTxt)
			t, _ := getTemplate("layouts/base.html", "user.html", "header.html")

			td := TmplData{currUser, User{
				Mail: email,
			}}
			t.Execute(w, td)
		} else {
			// render success template
			t, _ := getTemplate("layouts/base.html", "user.html", "header.html")
			td := TmplData{currUser, User{
				UID:  uid,
				Pass: pw,
			}}
			t.Execute(w, td)
		}
	}
}

func handleText(w http.ResponseWriter, r *http.Request) {
	currUser := r.Context().Value(currUserKey).(User)
	if r.Method == http.MethodGet {
		t, _ := getTemplate("layouts/base.html", "text.html", "header.html")
		td := TmplData{currUser, nil}
		t.Execute(w, td)
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Println(err)
		}
		text := r.PostForm.Get("text")
		if len(text) < 200 {
			log.Println("Text provided is to short for the purpose of this experiment")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		fre := calculateFRE(text)
		log.Println(fre)
		if fre <= 70. {
			log.Println("Text provided is to hard for the purpose of this experiment")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		txt := &Text{
			TID:     primitive.NewObjectID(),
			Content: text,
			FRE:     fre,
		}
		addNewText(txt, "../../web/data/texts", currUser.UID)
		w.WriteHeader(http.StatusOK)
	} else if r.Method == http.MethodPatch {
		scanDirForNewTexts("../../web/data/texts")
	}
}

func handleFRE(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Println(err)
		}
		txt := r.PostForm.Get("text")
		fre := calculateFRE(txt)
		fmt.Fprintf(w, "%f", fre)
	} else if r.Method == http.MethodPatch {
		updateTextFREs("../../web/data/texts")
		scanDirForNewTexts("../../web/data/texts")
	}
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	currUser := r.Context().Value(currUserKey).(User)
	if r.Method == http.MethodGet {
		tttid := r.URL.Query().Get("tttid")

		res := getTest(tttid, &currUser)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(res)
		if err != nil {
			log.Println(err)
		}
	} else if r.Method == http.MethodPost {
		tt := &TypingTest{}
		err := json.NewDecoder(r.Body).Decode(tt)
		if err != nil {
			log.Println(err)
		}
		if tt.UID != currUser.UID {
			log.Println("There was an attempt to save a test for a different user")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		err = saveTypingTestResult(tt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		currUser.CurrSession.LastTTID = tt.TTID
		currUser.CurrSession.LastKID = tt.KID
		updateSession(&currUser)
	} else if r.Method == http.MethodPatch {
		// scanDirForNewTexts("../../web/data/texts")
		// generateTypingTestFiles()
		updateTypingTestTextFREs("../../web/data/typing_test_texts")
		scanDirForNewTypingTestTexts("../../web/data/typing_test_texts")
	}
}

func handleAdminPanel(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := getTemplate("layouts/base.html", "admin_panel.html", "header.html")
		currUser := r.Context().Value(currUserKey).(User)
		td := TmplData{currUser, nil}
		t.Execute(w, td)
	}
}

func getTemplate(files ...string) (*template.Template, error) {
	tmplFolder := "../../web/templates/"
	for i, f := range files {
		files[i] = tmplFolder + f
	}
	// Name has to be the basename of at least one of the template files
	t := template.New(path.Base(files[0]))
	// Add custom helper funcs to template
	t.Funcs(template.FuncMap{
		"UnwrapOID":         unwrapObjectID,
		"GetSurveys":        getSurveys,
		"GetAvailableChars": getAvailableChars,
	})

	t, err := t.ParseFiles(files...)

	return t, err
}
