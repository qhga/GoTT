/*
  TT: Typing test, written for my bachelor thesis
  author: phga
  date: 2020-10-23
*/
// Definitions
let letters   = [];   // Array for ready-to-display letters of P
let position  = 0;    // Current position in letters
let timer     = null; // Reference for the timer used to time the individual typing tests
let TEST_TIME = 0;    // Time per test: Defined in the config.json file
let timeLeft  = 0;    // Time left for active test

// Variables for measurements
let P   = []; // Presented Text: Text the participant has to transcribe
let T   = []; // Transcript: Text written by the participant
let IS  = ""; // Input Stream: Anything written by the participant
let TL  = 0;  // Transcript Length: Correct and Incorrect input
let ISL = 0;  // Input Stream Length: Any Input during transcription
let C   = 0;  // Correct: Correct Input
let INF = 0;  // Incorrect Not Fixed: Remaining IF
let F   = 0;  // Fixes: Keystrokes to fix incorrect input
let IF  = 0;  // Inncorrect Fixed: All characters backspaced during input
let letter_errors = {};

// Results for error measurements
let CER  = 0;  // Corrected Error Rate: IF/(C + INF + IF)
let UER  = 0;  // Uncorrected Error Rate: INF/(C + INF + IF)
let TER  = 0;  // Total Error Rate: (INF + IF)/(C + INF + IF)
let KSPC = 0;  // Keystrokes per Character: ISL/TL
let KSPS = 0;  // Keystrokes per Second: ISL - 1 / TEST_TIME [s]

// Results for entry performance
let a        = 1; // Penalty exponent: [1;inf[, usually 1
let wpm      = 0; // Words Per Minute: (TL/5)/time
let cwpm     = 0; // Corrected WPM: wpm * (1 - UER)^a
let accuracy = 0; // Accuracy: TL - (IF + INF) / TL * 100

// Typing test object later send to the backend
let typingTest = {};

// Scrolling parameters
let lineHeight = 25;
let scrollHeight = 0;
let scrolls = 0;
let SCROLL_THRESH = 5;

// Timer parameters
let fgc = "#fafafa";
let bgc = "#f1f2f6";

// Initialisations
let allowedChars = /^.$|Backspace|Shift/; // Allowed chars for the test
let ttArea = document.getElementById("tt-area"); // Main test area
let startButton = document.getElementById("tt-start-button");
let resetButton = document.getElementById("tt-reset-button");
let timerBar = document.getElementById("tt-timer");
let infoBox = document.getElementById("tt-info");
let result = document.getElementById("tt-result");
let keyboardSelector = document.getElementById("tt-keyboard-selector");
let ttSelector = document.getElementById("tt-typing-test-selector");

// Init the Test
startButton.addEventListener("click", async e => {
    e.preventDefault();
    clearTimer();
    // wait for the response from the server
    typingTest = await fetchNewTypingTest();
    TEST_TIME = typingTest.test_time;
    // Initialize the test area with text
    initTTArea([...typingTest.test_text]);
    resetTestParams();
    setInfo(infoBox, "Timer begins as soon as you start typing", "black", "grey1");
    toggleHidden(infoBox, "off");
    toggleHidden(result, "on");
    toggleDisabled(keyboardSelector, "off");
    toggleDisabled(ttSelector, "off");
    toggleHidden(keyboardSelector, "off");
    // Typing events
    document.addEventListener("keydown", testRoutine);
})

// Prevent triggering start button with space
startButton.addEventListener("keyup", async e => {
    e.preventDefault();
});

// Prevent triggering start button with enter
startButton.addEventListener("keypress", async e => {
    e.preventDefault();
});


resetButton.addEventListener("click", e => {
    e.preventDefault();
    document.removeEventListener("keydown", testRoutine);
    clearTimer();
    initTTArea([""]);
    toggleHidden(infoBox, "on");
    toggleHidden(result, "on");
    toggleDisabled(keyboardSelector, "off");
    toggleDisabled(ttSelector, "off");
    toggleHidden(keyboardSelector, "off");
    resetTestParams();
})

function setTimer(t) {
    timeLeft = t; // Initialize test time
    toggleHidden(infoBox, "on");
    toggleDisabled(keyboardSelector, "on");
    toggleDisabled(ttSelector, "on");
    timer = setInterval(() => {
        timeLeft--;

        // If in config file show_timer is true
        if (typingTest.test_show_timer){
            let p = timeLeft / TEST_TIME * 100;
            timerBar.style.background = `linear-gradient(90deg, ${bgc} 0%, ${bgc} ${p}%, ${fgc} ${p}%)`;
        }

        if (timeLeft == 0) {
            document.removeEventListener("keydown", testRoutine);
            clearTimer();
            toggleDisabled(keyboardSelector, "off");
            toggleHidden(keyboardSelector, "on");
            saveResults();

            // If in config file show_result is true
            if (typingTest.test_show_result) {
                showResult();
            }
        }
    }, 1000);
}

function clearTimer() {
    clearInterval(timer);
    timer = null;
    timerBar.style.background = bgc;
}

function checkInput(l) {
    let cl = letters[position];

    if (l == "Backspace") {
        if (position == 0) {
            F++;
            return;
        }
        // Increase Fixes and Incorrect Fixes
        F++;
        IF++;

        cl = letters[--position];
        // Add Backspace char to InputStream
        IS += "<";
        // Remove last Char from T
        T.pop();
        // Fixed an Error
        if (cl.classList.contains("tt-error")){
            INF--; // Incorrect not Fixed
        }
        cl.classList.remove("tt-error", "tt-done");
        cl.classList.add("tt-todo");

        return;
    } else if (l == cl.textContent) {
        // correct input => mark done, forward one letter
        cl.classList.remove("tt-todo");
        cl.classList.add("tt-done");
        T[position] = l;
        // Add Backspace char to InputStream
        IS += l;
    } else {
        // incorrect input => mark error, forward one letter
        cl.classList.remove("tt-todo");
        cl.classList.add("tt-error");
        T[position] = l;
        INF++;
        // Add erroneous char to InputStream, prefix with _
        IS += `_${l}`;
        // Count errors per letter
        if (letter_errors[cl.textContent] == undefined) {
            letter_errors[cl.textContent] = 1;
        } else {
            letter_errors[cl.textContent]++;
        }
    }
    if (cl.classList.contains("tt-newline")) {
        scrollDown();
        cl.classList.remove("tt-newline");
    }
    position++;
}

function initTTArea(rawLetters) {
    // Reset scroll and cursor position
    resetTTArea();
    position = 0;

    // Initialize P
    P = rawLetters;

    // Prepare Text for presentation
    let ttAreaHtml = "";
    rawLetters.forEach(l => {
        ttAreaHtml += `<span class="tt-todo">${l}</span>`;
    });
    ttArea.innerHTML = ttAreaHtml;
    letters = ttArea.childNodes;
    // Mark the letters before line-wraps
    markWrapLetters();
}

function testRoutine(e) {
    e.preventDefault();
    if (null == timer) {
        // Start timer
        setTimer(TEST_TIME);
    }

    let l = e.key;

    if (allowedChars.test(l)) {
        checkInput(l);
    }
}

function markWrapLetters() {
    let top = letters[0].getBoundingClientRect().top;
    let nTop = 0; // New Top value
    let prevL = letters[0]; // Letter from prev Iter

    // Iterate through letter DOM objects
    letters.forEach(l => {
        nTop = l.getBoundingClientRect().top;
        if (nTop > top) {
            top = nTop;
            prevL.classList.add("tt-newline");
        }
        prevL = l;
    });
}

function scrollDown() {
    scrolls++;
    if (scrolls > SCROLL_THRESH) {
        scrollHeight -= lineHeight;
        ttArea.style.marginTop = scrollHeight + "px";
    }
}

function resetTestParams() {
    // Variables for measurements
    P   = []; // Presented Text: Text the participant has to transcribe
    T   = []; // Transcript: Text written by the participant
    IS  = ""; // Input Stream: Anything written by the participant
    TL  = 0;  // Transcript Length: Correct and Incorrect input
    ISL = 0;  // Input Stream Length: Any Input during transcription
    C   = 0;  // Correct: Correct Input
    INF = 0;  // Incorrect Not Fixed: Remaining IF
    F   = 0;  // Fixes: Keystrokes to fix incorrect input
    IF  = 0;  // Inncorrect Fixed: All characters backspaced during input
    letter_errors = {};

    // Results for error measurements
    CER  = 0;  // Corrected Error Rate: IF/(C + INF + IF)
    UER  = 0;  // Uncorrected Error Rate: INF/(C + INF + IF)
    TER  = 0;  // Total Error Rate: (INF + IF)/(C + INF + IF)
    KSPC = 0;  // Keystrokes per Character: ISL/TL
    KSPS = 0;  // Keystrokes per Second: ISL - 1 / TEST_TIME [s]

    // Results for entry performance
    wpm      = 0; // Words Per Minute: (TL/5)/time
    cwpm     = 0; // Corrected WPM: wpm * (1 - UER)^a
    accuracy = 0; // Accuracy: TL - (IF + INF) / TL * 100
}

function resetTTArea() {
    scrolls = 0;
    scrollHeight = 0;
    ttArea.style.marginTop = "0px";
}

function roundToPrecision(num, p) {
    x = Math.pow(10, p);
    return Math.round(num * x) / x;
}

function calculateResult() {
    // Required values for calculation
    mins = TEST_TIME / 60;
    TL = T.length; // == C + INF
    ISL = TL + F + IF;
    C  = TL - INF;

    CER = roundToPrecision(IF / (TL + IF), 5); // OK
    UER = roundToPrecision(INF / (TL + IF), 5); // OK
    TER = roundToPrecision((INF + IF)/(TL + IF), 5); // OK
    KSPC = roundToPrecision(ISL / TL, 5); // OK
    KSPS = roundToPrecision((ISL - 1) / TEST_TIME, 5); // OK

    // TL - 1 because the first char is entered at 0 seconds
    wpm = roundToPrecision((TL - 1) / (5 * mins), 2); // OK
    cwpm = roundToPrecision(wpm * Math.pow((1 - UER), a), 2); // OK
    accuracy = roundToPrecision(C / (TL + IF) * 100, 2);

    wpm = (wpm > 0) ? wpm : 0;
    cwpm = (cwpm > 0) ? cwpm : 0;
    accuracy = (accuracy > 0) ? accuracy : 0;
}

function showResult() {
    result.innerHTML = `<b>WPM:</b> ${wpm} - <b>CWPM:</b> ${cwpm} - <b>ACC:</b> ${accuracy}% - <b>CER:</b> ${CER} - <b>UER:</b> ${UER} - <b>TER:</b> ${TER} - <b>IF:</b> ${IF} - <b>INF:</b> ${INF} - <b>F:</b> ${F} - <b>TL:</b> ${TL} - <b>KSPC:</b> ${KSPC} - <b>KSPS:</b> ${KSPS}`;
    toggleHidden(result, "off");
}

function hideResult() {
    result.style.visibility = `hidden`;
}

function saveResults() {
    calculateResult();
    typingTest.TL = TL;
    typingTest.ISL = ISL;
    typingTest.INF = INF;
    typingTest.IF = IF;
    typingTest.F = F;
    typingTest.WPM = wpm;
    typingTest.CWPM = cwpm;
    typingTest.CER = CER;
    typingTest.UER = UER;
    typingTest.TER = TER;
    typingTest.KSPC = KSPC;
    typingTest.KSPS = KSPS;
    typingTest.accuracy = accuracy;
    typingTest.test_text = "";
    typingTest.kid = keyboardSelector.options[keyboardSelector.selectedIndex].value;
    typingTest.letter_errors = letter_errors;
    typingTest.IS = IS;
    postTypingTestResult(`/test`);
}

async function fetchNewTypingTest() {
    let reply = await fetch(`/test?tttid=${ttSelector.options[ttSelector.selectedIndex].value}`);
    let test = reply.json();
    return test;
}

function postTypingTestResult(url) {
    fetch(url, {
        method: "POST",
        body: JSON.stringify(typingTest),
        headers: {
            "Content-Type": "application/json; charset=UTF-8"
        }
    }).then(res => {
        if (res.status == 200) {
            setInfo(infoBox, "Result saved successfully!", "green", "black");
            toggleHidden(infoBox, "off");
        } else {
            setInfo(infoBox, "There was an error while saving your result!", "red", "grey1");
            toggleHidden(infoBox, "off");
        }
        setTimeout(() => {
            toggleHidden(infoBox, "on");
            setInfo(infoBox, "Timer begins as soon as you start typing", "black", "grey1");
        }, 5000);
    });

}