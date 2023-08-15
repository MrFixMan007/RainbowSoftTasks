"use strict";
var JSONValue = (function () {
    function JSONValue(Type, Name, Size) {
        this.Type = Type;
        this.Name = Name;
        this.Size = Size;
    }
    return JSONValue;
}());
var defaultRoot = '/';
var backButton = document.getElementById('backButton');
var sortType = document.getElementById('sortType');
if (backButton) {
    backButton.addEventListener('click', getBackDir);
}
if (sortType) {
    sortType.addEventListener('click', getDir);
}
var spinnerLoadDir = document.createElement('div');
spinnerLoadDir.className = "loader";
document.body.append(spinnerLoadDir);
var root = document.getElementById('root');
if (root)
    root.innerHTML = defaultRoot;
getDir();
function getBackDir() {
    var rootStr = "/";
    if (root)
        rootStr = String(root.textContent);
    var lastIndexOfSlesh = rootStr.lastIndexOf("/");
    if (root)
        root.textContent = rootStr.slice(0, lastIndexOfSlesh);
    if (root.textContent == "/" || root.textContent == "") {
        root.innerHTML = defaultRoot;
    }
    getDir();
}
function getDir() {
    if (spinnerLoadDir)
        spinnerLoadDir.classList.remove('hidden');
    var seconds = 0;
    var timer = setInterval(function () {
        seconds++;
    }, 10);
    var sortType = document.getElementById('sortType');
    var xhr = new XMLHttpRequest();
    var url = new URL("http://".concat(window.location.host, "/dir"));
    var roorStr;
    if (root) {
        roorStr = String(root.textContent);
        url.searchParams.set('root', roorStr);
    }
    if (sortType.checked == true) {
        url.searchParams.set('sortType', 'asc');
    }
    else
        url.searchParams.set('sortType', 'desc');
    xhr.open('GET', url);
    xhr.send();
    xhr.onload = function () {
        renderDir(xhr);
        var divTimer = document.getElementById('timer');
        clearInterval(timer);
        divTimer.innerHTML = "\u0412\u0440\u0435\u043C\u044F \u0432\u044B\u043F\u043E\u043B\u043D\u0435\u043D\u0438\u044F: ".concat(seconds / 100, " \u0441\u0435\u043A\u0443\u043D\u0434(\u044B)");
    };
    xhr.onerror = function () {
        alert('[Ошибка соединения]');
    };
}
;
function renderDir(xhr) {
    spinnerLoadDir.classList.add('hidden');
    var divUnswers = document.getElementById('unswers');
    divUnswers.innerHTML = "";
    var unmarshFiles = JSON.parse(xhr.response);
    if (unmarshFiles == null) {
        divUnswers.innerText = "Папка пуста";
        return;
    }
    var ul = document.createElement('ul');
    ul.className = "filesUl";
    ul.id = "files";
    divUnswers.appendChild(ul);
    for (var i = 0; i < unmarshFiles.length; i++) {
        var li = document.createElement("li");
        li.id = "filesLi".concat(i);
        if (unmarshFiles[i].Type == "file")
            li.className = "lis fileLi";
        else
            li.className = "lis folderLi";
        var rootString = String(root === null || root === void 0 ? void 0 : root.textContent);
        li.innerHTML = "".concat(unmarshFiles[i].Name.slice(rootString.length + 1), "&nbsp;:&nbsp").concat(unmarshFiles[i].Size, " \u0431\u0430\u0439\u0442(\u043E\u0432)");
        ul.onclick = function (event) {
            var eventTarget = event.target;
            var eventTargetString = String(event.target);
            if (eventTargetString == '[object HTMLUListElement]')
                return;
            var dots = document.getElementsByClassName('lis');
            var clickedFileLi = unmarshFiles[Array.from(dots).indexOf(eventTarget)];
            if (clickedFileLi.Type == 'file')
                return;
            root.innerHTML = clickedFileLi.Name;
            getDir();
        };
        ul.append(li);
    }
}
