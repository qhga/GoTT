package main

import (
	"testing"
)

func TestCountSyllables(t *testing.T) {
	txt := "Schwimmen"
	got := countSyllables(txt)
	want := 2
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}
	txt = "Knacken"
	got = countSyllables(txt)
	want = 2
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}
	txt = "Magnet"
	got = countSyllables(txt)
	want = 2
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}
	txt = "Kochen"
	got = countSyllables(txt)
	want = 2
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}
	txt = "Polizeiwachstation"
	got = countSyllables(txt)
	want = 7
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}
	txt = "Versicherungsmaklerinnen sind supertolle Weihnachtstrolle."
	got = countSyllables(txt)
	want = 17
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}
	txt = "Ich ich ich. Ich ICH ICH."
	got = countSyllables(txt)
	want = 6
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}
	txt = ""
	got = countSyllables(txt)
	want = 0
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}
	txt = "a l l e s K l a r o P o l i z e i"
	got = countSyllables(txt)
	want = 0
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}
	txt = "Ich bin mir insgesamt nicht ehrlich sicher."
	got = countSyllables(txt)
	want = 11
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}
	txt = "Das tut mir aua machen!"
	got = countSyllables(txt)
	want = 6
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}
	txt = "Hallo mein Name ist Bernd und ich bin ein Brot. Mir geht es super mit allem denn ich bin eine Gurke. Nein aua das tut mir weh!"
	got = countSyllables(txt)
	want = 33
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}

	txt = "Es war nicht einfach, aber sie kämpfte darum, mit der neuen Situation gut zurechtzukommen. Vor allem die extreme Höhe machte ihr zu schaffen. In der ersten Nacht am Berg hatte sie kein Auge zugemacht."
	got = countSyllables(txt)
	want = 58
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}

	// Ne Si tu ti
	txt = "Neuen Situation aber. Nationen Brauerrei dualist duelist"
	got = countSyllables(txt)
	want = 22
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}

	txt = "Höhe kämpfte"
	got = countSyllables(txt)
	want = 4
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}

	// Ka se ra Rä ub ei ni ro le Wi se
	// Ka se ra uu Rä ub ei ni ro le Wi se
	txt = "Kaiser Trauung Räuber einig Broiler Wiese"
	got = countSyllables(txt)
	want = 12
	if got != want {
		t.Errorf("Syllable count incorrect, got: %d, want: %d", got, want)
	}
}

func TestCountWords(t *testing.T) {
	txt := "Hallo ich bin der Philip. Alles ist gut!"
	got := countWords(txt)
	want := 8
	if got != want {
		t.Errorf("Word count incorrect, got: %d, want: %d", got, want)
	}
	txt = "h a l l o"
	got = countWords(txt)
	want = 0
	if got != want {
		t.Errorf("Word count incorrect, got: %d, want: %d", got, want)
	}
	txt = ""
	got = countWords(txt)
	want = 0
	if got != want {
		t.Errorf("Word count incorrect, got: %d, want: %d", got, want)
	}
	txt = "In der Nacht!"
	got = countWords(txt)
	want = 3
	if got != want {
		t.Errorf("Word count incorrect, got: %d, want: %d", got, want)
	}
	txt = "Hallo mein Name ist Bernd und ich bin ein Brot. Mir geht es super mit allem denn ich bin eine Gurke. Nein aua das tut mir weh!"
	got = countWords(txt)
	want = 27
	if got != want {
		t.Errorf("Word count incorrect, got: %d, want: %d", got, want)
	}

	txt = "örtörtört ärtätött ßttßtt süusüsss"
	got = countWords(txt)
	want = 4
	if got != want {
		t.Errorf("Word count incorrect, got: %d, want: %d", got, want)
	}
}

func TestCountSentences(t *testing.T) {
	txt := "Hallo ich bin der Philip. Alles ist gut!"
	got := countSentences(txt)
	want := 2
	if got != want {
		t.Errorf("Sentence count incorrect, got: %d, want: %d", got, want)
	}
	txt = "h a l l o"
	got = countSentences(txt)
	want = 0
	if got != want {
		t.Errorf("Sentence count incorrect, got: %d, want: %d", got, want)
	}
	txt = ""
	got = countSentences(txt)
	want = 0
	if got != want {
		t.Errorf("Sentence count incorrect, got: %d, want: %d", got, want)
	}
	txt = "In der Nacht! Ich. Bin nicht? Der; Den du kennst."
	got = countSentences(txt)
	want = 5
	if got != want {
		t.Errorf("Sentence count incorrect, got: %d, want: %d", got, want)
	}
	txt = "In der Nacht!"
	got = countSentences(txt)
	want = 1
	if got != want {
		t.Errorf("Sentence count incorrect, got: %d, want: %d", got, want)
	}
	txt = "Hallo mein Name ist Bernd und ich bin ein Brot. Mir geht es super mit allem denn ich bin eine Gurke. Nein aua das tut mir weh!"
	got = countSentences(txt)
	want = 3
	if got != want {
		t.Errorf("Sentence count incorrect, got: %d, want: %d", got, want)
	}
}

func TestCalculateFRE(t *testing.T) {
	txt := "Hallo mein Name ist Bernd und ich bin ein Brot. Mir geht es super mit allem denn ich bin eine Gurke. Nein aua das tut mir weh!"
	got := calculateFRE(txt)
	want := 99.5
	if got != want {
		t.Errorf("FRE incorrect, got: %f, want: %f", got, want)
	}
}
