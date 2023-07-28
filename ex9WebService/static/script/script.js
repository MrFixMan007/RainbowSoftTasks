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
  var seconds = 0;
  const timer = setInterval(()=>
  {
    seconds++;
  }, 10);

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
    renderDir(xhr);
    const divTimer = document.getElementById('timer');
    clearInterval(timer)
    divTimer.innerHTML=`Время выполнения: ${seconds/100} секунд(ы)`;
  };

  xhr.onerror = function() { 
    alert('[Ошибка соединения]');
  };
};

function renderDir(xhr){
  const divUnswers = document.getElementById('unswers');
  divUnswers.innerHTML = "";

  unmarshFiles = JSON.parse(xhr.response);
  if(unmarshFiles == null){
    divUnswers.innerText = "Папка пуста";
    return;
  }

  let ul = document.createElement('ul');
  ul.className = "filesUl";
  ul.id = "files";
  divUnswers.appendChild(ul);

  for (let i = 0; i < unmarshFiles.length; i++){
    var li = document.createElement("li");
    li.id = `filesLi${i}`;
    if(unmarshFiles[i].Type == "file") li.className = "lis fileLi"
    else li.className = "lis folderLi"
    li.innerHTML = `${unmarshFiles[i].Name.slice(root.length + 1)}&nbsp;:&nbsp${unmarshFiles[i].Size} байт(ов)`;

    ul.onclick = (event) => {
      var dots = document.getElementsByClassName('lis'); //получаем все li
      clickedFileLi = unmarshFiles[Array.from(dots).indexOf(event.target)] //получаем json-объект по индексу, найденному в массиве li через event.target

      if(clickedFileLi.Type == 'file') return
      document.getElementById('root').value = clickedFileLi.Name; //меняем строку ввода, если папка
      getDir()
    }
    ul.append(li);
  }
}