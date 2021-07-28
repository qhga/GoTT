function toggleDisabled(obj, state) {
    if (state == "on") {
        obj.setAttribute("disabled", "");
    } else {
        obj.removeAttribute("disabled");
    }
}

function toggleHidden(obj, state) {
    if (state == "off") {
        obj.classList.remove("hidden");
    } else {
        obj.classList.add("hidden");
    }
}

function setInfo(obj, txt, bcolor, fcolor) {
    obj.innerText = txt;
    obj.style.background = `var(--${bcolor})`;
    obj.style.color = `var(--${fcolor})`;
}