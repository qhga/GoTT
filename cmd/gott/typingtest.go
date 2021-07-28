package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var textCollection []Text
var typingTestTextCollection map[string]TypingTestText
var typingTestTextTTIDs []string

// Text snippet used in the TypingTest
type Text struct {
	TID     primitive.ObjectID `json:"tid"`
	Content string             `json:"content"`
	FRE     float64            `json:"fre"`
}

type TypingTestText struct {
	TTTID    string               `json:"tttid" bson:"tttid"`
	TIDs     []primitive.ObjectID `json:"tids" bson:"tids"`
	MeanFRE  float64              `json:"mean_fre" bson:"mean_fre"`
	TestText string               `json:"test_text" bson:"-"`
}

// Time is not required and can be extracted from TTID (First 4 bytes)
// Explanation for the individual fields can be found in tt.js
type TypingTest struct {
	TTID    primitive.ObjectID   `json:"ttid" bson:"_id"`
	UID     string               `json:"uid" bson:"uid"`
	TTTID   string               `json:"tttid" bson:"tttid"`
	TIDs    []primitive.ObjectID `json:"tids" bson:"tids"`
	KID     string               `json:"kid" bson:"kid"`
	MeanFRE float64              `json:"mean_fre" bson:"mean_fre"`

	TL             int     `json:"TL" bson:"TL"`
	ISL            int     `json:"ISL" bson:"ISL"`
	INF            int     `json:"INF" bson:"INF"`
	IF             int     `json:"IF" bson:"IF"`
	F              int     `json:"F" bson:"F"`
	WPM            float64 `json:"WPM" bson:"WPM"`
	CWPM           float64 `json:"CWPM" bson:"CWPM"`
	CER            float64 `json:"CER" bson:"CER"`
	UER            float64 `json:"UER" bson:"UER"`
	TER            float64 `json:"TER" bson:"TER"`
	KSPC           float64 `json:"KSPC" bson:"KSPC"`
	KSPS           float64 `json:"KSPS" bson:"KSPS"`
	Accuracy       float64 `json:"accuracy" bson:"accuracy"`
	TestText       string  `json:"test_text" bson:"-"`
	TestTime       int     `json:"test_time" bson:"test_time"`
	TestShowTimer  bool    `json:"test_show_timer" bson:"-"`
	TestShowResult bool    `json:"test_show_result" bson:"-"`

	LetterErrors map[string]int `json:"letter_errors" bson:"letter_errors"`
	IS           string         `json:"IS" bson:"IS"`
}

func init() {
	scanDirForNewTexts("../../web/data/texts")
	scanDirForNewTypingTestTexts("../../web/data/typing_test_texts")
	typingTestTextTTIDs = make([]string, len(typingTestTextCollection))

	// Gather TTIDs
	i := 0
	for key := range typingTestTextCollection {
		typingTestTextTTIDs[i] = key
		i++
	}

	sort.Strings(typingTestTextTTIDs)
}

func updateTypingTestTextFREs(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("Error while scanning text directory: ", err)
		return
	}

	for _, f := range files {
		if !f.IsDir() {
			content, err := ioutil.ReadFile(dir + "/" + f.Name())
			var text TypingTestText
			err = json.Unmarshal(content, &text)
			if err != nil {
				log.Println("Error while unmarshalling json: ", err)
			}
			text.MeanFRE = calculateFRE(text.TestText)

			f, err := os.Create(dir + "/" + f.Name())
			defer f.Close()
			e := json.NewEncoder(f)
			e.SetIndent("", "    ")
			e.Encode(text)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func updateTextFREs(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("Error while scanning text directory: ", err)
		return
	}

	for _, f := range files {
		if !f.IsDir() {
			content, err := ioutil.ReadFile(dir + "/" + f.Name())
			var texts []Text
			err = json.Unmarshal(content, &texts)
			if err != nil {
				log.Println("Error while unmarshalling json: ", err)
			}
			for idx := range texts {
				texts[idx].FRE = calculateFRE(texts[idx].Content)
			}

			f, err := os.Create(dir + "/" + f.Name())
			defer f.Close()
			e := json.NewEncoder(f)
			e.SetIndent("", "    ")
			e.Encode(texts)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// scanDirForNewTexts scans the given path for json files and unmarshals the
// json into Text Objects which are then stored in the global var textCollection
func scanDirForNewTexts(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("Error while scanning text directory: ", err)
		return
	}

	// Reset textCollection
	// (0 Is Important, otherwise there will be one empty entry)
	textCollection = make([]Text, 0)

	for _, f := range files {
		if !f.IsDir() {
			content, err := ioutil.ReadFile(dir + "/" + f.Name())
			if err != nil {
				log.Println("Error while reading file: ", err)
			}

			var texts []Text
			err = json.Unmarshal(content, &texts)
			if err != nil {
				log.Println("Error while unmarshalling json: ", err)
				continue
			}
			textCollection = append(textCollection, texts...)
		}
	}
}

func scanDirForNewTypingTestTexts(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("Error while scanning TypingTestText directory: ", err)
		return
	}

	typingTestTextCollection = make(map[string]TypingTestText)

	for _, f := range files {
		if !f.IsDir() {
			content, err := ioutil.ReadFile(dir + "/" + f.Name())
			if err != nil {
				log.Println("Error while traversing directory: ", err)
			}

			var typingTestText TypingTestText
			err = json.Unmarshal(content, &typingTestText)
			if err != nil {
				log.Println("Error while unmarshalling json: ", err)
				continue
			}
			typingTestTextCollection[typingTestText.TTTID] = typingTestText
		}
	}
}

func saveTypingTestResult(tt *TypingTest) error {
	_, err := colTypingTests.InsertOne(mongoCtx, tt)
	if err != nil {
		log.Println(err)
	}
	return err
}

func getTypingTestResults(u *User) (*[]TypingTest, error) {
	filter := bson.D{}
	if u != nil {
		filter = bson.D{{"uid", u.UID}}
	}

	cursor, err := colTypingTests.Find(mongoCtx, filter)
	if err != nil {
		log.Fatal("Could not find any results in the DB: ", err)
	}

	var tts []TypingTest
	err = cursor.All(mongoCtx, &tts)
	if err != nil {
		log.Fatal("Something went wrong while gathering the results: ", err)
	}

	return &tts, err
}

func (t *Text) getLength() int {
	return len(t.Content)
}

func getAvailableChars() (charCount int) {
	scanDirForNewTexts("../../web/data/texts")
	for _, text := range textCollection {
		charCount += text.getLength()
	}
	return
}

func getTest(tttid string, currUser *User) (res []byte) {
	var tt *TypingTest
	if tttid == "random" {
		scanDirForNewTexts("../../web/data/texts")
		tt = generateNewTypingTest(currUser, 20)
		// handle index out of bound
	} else {
		ttt := typingTestTextCollection[tttid]
		tt = &TypingTest{
			TTID:           primitive.NewObjectID(),
			UID:            currUser.UID,
			TTTID:          tttid,
			MeanFRE:        ttt.MeanFRE,
			TestText:       ttt.TestText,
			TestTime:       config.TestTime,
			TestShowTimer:  config.TestShowTimer,
			TestShowResult: config.TestShowResult,
		}
	}
	var err error
	res, err = json.Marshal(tt)
	if err != nil {
		log.Println("Generating text snippets failed: ", err)
		return
	}
	return
}

func generateTypingTestFiles() {
	// Generate 10 Typing Tests with enough text for fast typists (160 wpm)
	// min_char_count := 10 * 160 * 5 * 5
	minutes := int(config.TestTime / 60)
	if minutes < 1 {
		minutes = 1
	}
	minCharCount := config.TestAmount * config.TestMaxWpm * 5 * minutes
	minCharPerTest := minCharCount / config.TestAmount
	// 40.000+ symbols are required for 10 typing tests
	availableChars := getAvailableChars()

	log.Println(minCharCount)
	log.Println(availableChars)

	if minCharCount > availableChars {
		log.Println("Not enough text snippets to generate Typing Tests")
		return
	}

	// Sort array according to FRE
	sort.Slice(textCollection, func(i, j int) bool {
		return textCollection[i].FRE < textCollection[j].FRE
	})

	one_or_two := 1

	for currText := 0; currText < config.TestAmount; currText++ {
		textIDs := make([]primitive.ObjectID, 0, minCharPerTest/100)
		testText := ""
		i := currText

		for currChars := 0; currChars < minCharPerTest; {
			textIDs = append(textIDs, textCollection[i].TID)
			currChars += textCollection[i].getLength()
			if i == currText {
				testText += textCollection[i].Content
			} else {
				testText += " " + textCollection[i].Content
			}
			i = (i + config.TestAmount) % len(textCollection)
		}

		ttt := &TypingTestText{
			TTTID:    fmt.Sprintf("T%d_%d", int(currText/2), one_or_two),
			TIDs:     textIDs,
			MeanFRE:  calculateFRE(testText),
			TestText: testText,
		}

		if one_or_two == 1 {
			one_or_two = 2
		} else {
			one_or_two = 1
		}

		addNewTypingTestText(ttt, "../../web/data/typing_test_texts")
	}
}

func generateNewTypingTest(u *User, nTexts int) *TypingTest {
	if len(textCollection) < nTexts {
		nTexts = len(textCollection)
	}
	shuffleTexts()
	textIDs := make([]primitive.ObjectID, nTexts)
	testText := ""

	for i := 0; i < nTexts; i++ {
		textIDs[i] = textCollection[i].TID
		if i == 0 {
			testText += textCollection[i].Content
		} else {
			testText += " " + textCollection[i].Content
		}
	}

	tt := &TypingTest{
		TTID:           primitive.NewObjectID(),
		UID:            u.UID,
		TIDs:           textIDs,
		TTTID:          "random",
		MeanFRE:        calculateFRE(testText),
		TestText:       testText,
		TestTime:       config.TestTime,
		TestShowTimer:  config.TestShowTimer,
		TestShowResult: config.TestShowResult,
	}
	return tt
}

func shuffleTexts() {
	rand.Shuffle(len(textCollection), func(i, j int) {
		textCollection[i], textCollection[j] = textCollection[j], textCollection[i]
	})
}

func adjustText(txt string) string {
	return strings.ReplaceAll(strings.ReplaceAll(txt, "\n", " "), "  ", " ")
}

func addNewText(txt *Text, dir string, uid string) {
	fName := dir + "/" + uid + ".json"
	f, err := os.Open(fName)
	txt.Content = adjustText(txt.Content)
	txts := []Text{}
	// If the file exists
	if err == nil {
		// Decode possible array of objects to structs
		d := json.NewDecoder(f)
		err = d.Decode(&txts)
		if err != nil {
			log.Println(err)
		}
	}

	// Append new text to old texts
	txts = append(txts, *txt)
	f, err = os.Create(fName)
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetIndent("", "    ")
	e.Encode(txts)
	if err != nil {
		log.Println(err)
	}
}

func addNewTypingTestText(ttt *TypingTestText, dir string) {
	fName := dir + "/" + ttt.TTTID + ".json"
	f, err := os.Open(fName)
	// Append new text to old texts
	f, err = os.Create(fName)
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetIndent("", "    ")
	e.Encode(ttt)
	if err != nil {
		log.Println(err)
	}
}

func countSyllables(txt string) int {
	// Either Con+Vov, Start of Line OR Whitespace Vov+Con, Vov+Vov+... (Monophthongen)
	rx := regexp.MustCompile(`(?i)[^aeiouäöüßy\W][aeiouäöüßy]|\b[aeiouäöüßy][^aeiouäöüßy\W]|\b[aeiouäöüy]{2,}|u[aeuo]|(on|er)\b|\B(a|o|u|e)\B`)
	extraConsonants := []string{"ck", "x", "ch", "x", "sch", "x", "st", "x", "gn", "x"}
	extraVowels := []string{"äu", "i", "ie", "i"}
	r := strings.NewReplacer(extraConsonants...)
	txt = r.Replace(txt)
	r = strings.NewReplacer(extraVowels...)
	txt = r.Replace(txt)
	// log.Println(rx.FindAllString(txt, -1))
	syllableCount := len(rx.FindAllStringIndex(txt, -1))
	return syllableCount
}

func countWords(txt string) int {
	rx := regexp.MustCompile(`[\wäöüß]{2,}`)
	return len(rx.FindAllStringIndex(txt, -1))
}

func countSentences(txt string) int {
	rx := regexp.MustCompile(`[\wäöüß]{2,}[\?\.!;]`)
	return len(rx.FindAllStringIndex(txt, -1))
}

// Flesch-Reading-Ease (German)
// FRE = 180 - ASL - (58.5 * ASW)
// ASL = Average Sentence Length = Words / Sentence
// ASW = Average Number of Syllables per Word = Syllables / Words
func calculateFRE(txt string) float64 {
	syc := countSyllables(txt)
	wc := countWords(txt)
	sec := countSentences(txt)

	asl := float64(wc) / float64(sec)
	asw := float64(syc) / float64(wc)

	fre := math.Round((180.-asl-(58.5*asw))*100) / 100 // OK

	// Max fre is 100 (<0 and >100 is allowed, though not relevant in this case)
	if fre > 100. {
		fre = 100.
	}
	if fre < 0. {
		fre = 0.
	}
	return fre
}
