myButton.addEventListener('click', getDir);
getDir()

function getDir(root = "/") {
    rootInInput = document.getElementById('root').value
    if(rootInInput != "" && rootInInput != null){
        root = rootInInput;
    }

    sortType = document.getElementById('sortType')

    let xhr = new XMLHttpRequest();
    let url = new URL(`http://${window.location.host}/dir`);
    
    url.searchParams.set('root', root);
    if(document.getElementById('sortType').checked == true){
      url.searchParams.set('sortType', 'asc');
    }
    else url.searchParams.set('sortType', 'desc');

    xhr.open('GET', url);
    xhr.send();
    i = 0;

    xhr.onload = function() {
      unmarshFiles = JSON.parse(xhr.response);

      const divUnswers = document.getElementById('unswers');
      divUnswers.innerHTML = "";

      for (let i = 0; i < unmarshFiles.length; i++){
        var divSrc = document.createElement("div");
        divSrc.id = 'unswer' + i;
        divSrc.className = "unswer"
        divUnswers.appendChild(divSrc);

        divSrc.textContent = unmarshFiles[i].Type + " " + unmarshFiles[i].Name + " " + unmarshFiles[i].Size;
      }
    };

    xhr.onerror = function() { 
      alert('[Ошибка соединения]');
    };
};
