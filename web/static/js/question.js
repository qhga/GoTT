let frm = document.getElementById("frm-add-question");
let optionContainer = document.getElementById("option-container");
let questionPreview = document.getElementById("question-preview");
let infoBox = document.getElementById("general-info");

let btnAddOption = document.createElement("button");
btnAddOption.id = "btn-o-add";
btnAddOption.innerText = "+";

let ph = "Undefined";

// holds options {label: "xyz"}
let question = {
    q_type: "1",
    q_text: "",
    q_answers: {},
    q_label_r: "",
    q_label_l: ""
}
let qAnswers = {
    radio: [ph],
    multi: [ph],
    ftext: [ph],
    range: [ph]
};

document.addEventListener("readystatechange", () => {
    redisplayOptions(question.q_type);
});

frm.addEventListener("click", e => {
    e.preventDefault();
    let ele = e.target;

    if (ele.classList.contains("btn-o-delete")){
        deleteOption(ele.dataset.index);
    } else if (ele.id == "btn-o-add") {
        addNewOption();
    } else if (ele.id == "btn-submit") {
        submitQuestion();
    }
})

frm.addEventListener("change", e => {
    let ele = e.target;

    if (ele.id == "sel-Qtype") {
        question.q_type = ele.options[ele.selectedIndex].value;
        redisplayOptions();
    } else if (ele.classList.contains("inp-o-label")) {
        changeOption(ele.dataset.index, ele.value);
    } else if (ele.id == "q-text") {
        question.q_text = ele.value;
        renderQuestionPreview();
    } else if (ele.id == "inp-q-labelL") {
        question.q_label_l = ele.value;
        renderQuestionPreview();
    } else if (ele.id == "inp-q-labelR") {
        question.q_label_r = ele.value;
        renderQuestionPreview();
    }
});

function redisplayOptions() {
    let currOptions = "";

    let ll = `<label for="inp-q-labelL">Left label for whole scale:</label>
<input id="inp-q-labelL" type="text" value="${question.q_label_l}"/>`
    let lr = `<label for="inp-q-labelL">Right label for whole scale:</label>
<input id="inp-q-labelR" type="text" value="${question.q_label_r}"/>`

    switch(question.q_type) {
    case "1":
        currOptions += ll;
        currOptions += lr;
        qAnswers.radio.forEach((o, i) => {
            currOptions += `<input class="inp-o-label" data-index="${i}" style="width: 80%" type="text" value="${o}"/><button class="btn-o-delete" data-index="${i}">-</button>`;
        });
        break;
    case "2":
        qAnswers.multi.forEach((o, i) => {
            currOptions += `<input class="inp-o-label" data-index="${i}" style="width: 80%" type="text" value="${o}"/><button class="btn-o-delete" data-index="${i}">-</button>`;
        });
        break;
    case "3":
        qAnswers.ftext.forEach((o, i) => {
            currOptions += `<input class="inp-o-label" data-index="${i}" style="width: 80%" type="text" value="${o}"/><button class="btn-o-delete" data-index="${i}">-</button>`;
        });
        break;
    case "4":
        currOptions += ll;
        currOptions += lr;
        qAnswers.range.forEach((o, i) => {
            currOptions += `<input class="inp-o-label" data-index="${i}" style="width: 80%" type="text" value="${o}"/><button class="btn-o-delete" data-index="${i}">-</button>`;
        });
        break;
    }
    optionContainer.innerHTML = currOptions;
    optionContainer.appendChild(btnAddOption);
    renderQuestionPreview();
}

// MAYB: Add above aswell -> simply another button above with shift
function addNewOption() {
    switch(question.q_type) {
    case "1":
        qAnswers.radio.push(ph)
        break;
    case "2":
        qAnswers.multi.push(ph)
        break;
    case "3":
        qAnswers.ftext.push(ph)
        break;
    case "4":
        qAnswers.range.push(ph)
        break;
    }
    redisplayOptions();
}

function deleteOption(i) {
    switch(question.q_type) {
    case "1":
        qAnswers.radio.splice(i, 1);
        break;
    case "2":
        qAnswers.multi.splice(i, 1);
        break;
    case "3":
        qAnswers.ftext.splice(i, 1);
        break;
    case "4":
        qAnswers.range.splice(i, 1);
        break;
    }
    redisplayOptions();
}

function changeOption(i, txt) {
    switch(question.q_type) {
    case "1":
        qAnswers.radio[i] = txt;
        break;
    case "2":
        qAnswers.multi[i] = txt;
        break;
    case "3":
        qAnswers.ftext[i] = txt;
        break;
    case "4":
        qAnswers.range[i] = txt;
        break;
    }
    redisplayOptions();
}

function renderQuestionPreview() {
    txt = `
                <h2>${question.q_text}</h2>
                <div class="row-survey">
`

    switch(question.q_type) {
    case "1":
        txt += `<div class="col-survey">${question.q_label_l}</div>`;
        qAnswers.radio.forEach((o, i) => {
            txt += `<div class="col-survey">
                      <label for="${i}">${o}</label>
                      <input id="${i}" name="XYZ" type="radio" value="${o}" required/>
                    </div>`
        });
        txt += `<div class="col-survey">${question.q_label_r}</div>`;
        break;
    case "2":
        qAnswers.multi.forEach((o, i) => {
            txt += `<div class="col-survey">
                      <label for="${i}">${o}</label>
                      <input id="${i}" name="XYZ${i}" type="checkbox" value="${o}" required/>
                    </div>`
        });
        break;
    case "3":
        qAnswers.ftext.forEach((o, i) => {
            txt += `<div class="col-survey">
                      <label for="${i}">${o}</label>
                      <input id="${i}" name="XYZ" type="text" value="" required/>
                    </div>`
        });
        break;
    case "4":
        txt += `<div class="col-survey">${question.q_label_l}</div>`;
        qAnswers.range.forEach((o, i) => {
            txt += `<div class="col-survey">
                      <label for="${i}">${o}</label>
                      <input id="${i}" name="XYZ" type="range" value="" required/>
                    </div>`
        });
        txt += `<div class="col-survey">${question.q_label_r}</div>`;
        break;

    }


    txt += ` </div> `;
    questionPreview.innerHTML = `<div id="frm-survey">${txt}</div>`;
}

function submitQuestion() {
    switch(question.q_type) {
    case "1":
        question.q_answers = qAnswers.radio;
        break;
    case "2":
        question.q_label_l = "";
        question.q_label_r = "";
        question.q_answers = qAnswers.multi;
        break;
    case "3":
        question.q_label_l = "";
        question.q_label_r = "";
        question.q_answers = qAnswers.ftext;
        break;
    case "4":
        question.q_answers = qAnswers.range;
        break;
    }
    fetch("/question", {
        method: "POST",
        body: JSON.stringify(question),
        headers: {
            "Content-Type": "application/json; charset=UTF-8"
        }
    }).then(res => {
             if (res.status == 200) {
                 setInfo(infoBox, "Success", "green", "black");
                 toggleHidden(infoBox, "off");
             } else {
                 setInfo(infoBox, "Failure", "red", "grey1");
                 toggleHidden(infoBox, "off");
             }
             setTimeout(() => {
                 toggleHidden(infoBox, "on");
             }, 2000);
    });
}