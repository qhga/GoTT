let textBox = document.getElementById("ta-content");
let fre = document.getElementById("fre");
let cc = document.getElementById("cc");
let frm = document.getElementById("frm-contribute-text");
let infoBox = document.getElementById("general-info");

textBox.addEventListener("keyup", updateInfos);

function updateInfos(e) {
    console.log(e.key);
    if (e.key == " " || e.key == "." || e.key == "!" || e.key == "?" || e.key == "Backspace") {
        fetchFRE("/fre");
        updateCharCount();
    }
}

async function fetchFRE(url) {
    let txt = textBox.value.trim() + ".";
    fetch(url, {
        method: "POST",
        body: "text="+txt,
        headers: {
            "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8"
        }
    }).then(res => res.text()).then(f => {
        f = Math.round(parseFloat(f) * 100) / 100;
        (!isNaN(f)) ? fre.innerText = f : console.log(f);
    });
}

function updateCharCount() {
    cc.innerText = textBox.value.length;
}

frm.addEventListener("submit", e => {
    e.preventDefault();
    fetch("/text", {
        method: "POST",
        body: "text="+textBox.value.trim(),
        headers: {
            "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8"
        }

    }).then(res => {
        if (res.status == 200) {
            setInfo(infoBox, "Success", "green", "black");
            toggleHidden(infoBox, "off");
            setTimeout(() => {
                toggleHidden(infoBox, "on");
                window.location = "/text";
            }, 1300);
        } else {
            setInfo(infoBox, "Failure", "red", "grey1");
            toggleHidden(infoBox, "off");
            setTimeout(() => {
                toggleHidden(infoBox, "on");
            }, 1500);
        }
    })
})