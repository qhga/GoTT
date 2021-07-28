package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	QtypeSingle       = "1"
	QtypeMulti        = "2"
	QtypeFreeText     = "3"
	QtypeAnalogSlider = "4"
)

var questionCollecion map[primitive.ObjectID]Question

type Question struct {
	QID      primitive.ObjectID `bson:"_id" json:"qid"`
	Qtype    string             `bson:"-" json:"q_type"`
	Qtext    string             `bson:"q_text" json:"q_text"`
	QlabelL  string             `bson:"-" json:"q_label_l"`
	QlabelR  string             `bson:"-" json:"q_label_r"`
	Qmin     string             `bson:"-" json:"q_min"`
	Qmax     string             `bson:"-" json:"q_max"`
	Qstep    string             `bson:"-" json:"q_step"`
	QAnswers []string           `bson:"-" json:"q_answers"`
	Answer   string             `bson:"answer" json:"answer"`
}

type Survey struct {
	ID        primitive.ObjectID `bson:"_id"`
	SID       string             `bson:"sid"`
	UID       string             `bson:"pid"`
	KID       string             `bson:"kid"`
	TTID      primitive.ObjectID `bson:"ttid"`
	Questions []Question         `bson:"questions"`
}

func init() {
	scanDirForNewQuestions("../../web/data/questions")
}

func scanDirForNewQuestions(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return
	}

	// (0 Is Important, otherwise there will be one empty entry)
	questionCollecion = make(map[primitive.ObjectID]Question, 0)

	for _, f := range files {
		if !f.IsDir() {
			qf, err := os.Open(dir + "/" + f.Name())
			if err != nil {
				log.Println(err)
			}

			q := &Question{}
			err = json.NewDecoder(qf).Decode(q)
			if err != nil {
				log.Println(err)
				continue
			}
			questionCollecion[q.QID] = *q
		}
	}
	// log.Printf("%+v", questionCollecion)
}

func (s *Survey) setQuestions(oids ...primitive.ObjectID) {
	for _, oid := range oids {
		s.Questions = append(s.Questions, questionCollecion[oid])
	}
}

func addNewQuestion(q *Question, dir string) {
	fName := dir + "/" + unwrapObjectID(q.QID) + ".json"
	f, err := os.Create(fName)
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetIndent("", "    ")
	e.Encode(q)
	if err != nil {
		log.Println(err)
	}
}

func getSurveys() map[string][]primitive.ObjectID {
	return config.Surveys
}

func saveSurvey(s *Survey) error {
	_, err := colSurveys.InsertOne(mongoCtx, s)
	if err != nil {
		log.Println(err)
	}
	return err
}
