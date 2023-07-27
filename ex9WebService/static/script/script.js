myButton.addEventListener('click', getDir);
backButton.addEventListener('click', getBackDir);
getDir()

function getBackDir(){
  rootInInput = document.getElementById('root').value
  rootInInput = String(rootInInput)
  if(rootInInput == "/"){
    return;
  }
  lastIndeOfSlesh = rootInInput.lastIndexOf("/")
  document.getElementById('root').value = rootInInput.slice(0, lastIndeOfSlesh)
  getDir()
}

function getDir() {
    rootDefault = "/"
    root = rootDefault
    rootInInput = document.getElementById('root').value
    if(rootInInput != "" && rootInInput != null){
        root = rootInInput;
    }
    document.getElementById('root').value = root;

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

    xhr.onload = function() {
      let ul = document.createElement('ul');
      ul.className = "filesUl";
      ul.id = "files";

      unmarshFiles = JSON.parse(xhr.response);

      const divUnswers = document.getElementById('unswers');
      divUnswers.innerHTML = "";
      divUnswers.appendChild(ul);

      for (let i = 0; i < unmarshFiles.length; i++){
        var li = document.createElement("li");
        li.id = i;
        if(unmarshFiles[i].Type == "file") li.className = "fileLi"
        else li.className = "folderLi"
        li.textContent = unmarshFiles[i].Name.slice(root.length + 1) + " " + unmarshFiles[i].Size;
        ul.onclick = (event) => {
          if(unmarshFiles[event.target.id].Type == 'file') return
          document.getElementById('root').value = unmarshFiles[event.target.id].Name;
          getDir()
        }
        ul.append(li);
      }
    };

    xhr.onerror = function() { 
      alert('[Ошибка соединения]');
    };
};