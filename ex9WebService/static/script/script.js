backButton.addEventListener('click', getBackDir)
sortType.addEventListener('click', getDir);

//пока данные грузятся, страницу перекрывает спинер
spinnerLoadDir = document.createElement('div')
spinnerLoadDir.className = "loader"
document.body.append(spinnerLoadDir)

document.getElementById('root').value = "/home" //начальное значение
getDir()

//getBackDir чистит строку вывода с конца до первого встречного слэша
//и вызывает getDir()
function getBackDir(){
  let root = document.getElementById('root').value
  root = String(root)
  if(root == "/"){
    return;
  }
  lastIndeOfSlesh = root.lastIndexOf("/")
  document.getElementById('root').value = root.slice(0, lastIndeOfSlesh)
  getDir()
}

//getDir отправляет запрос серверу, обрабатывает ответы, замеряет время выполнения
function getDir() {  

  //ставим спинер, скрывая страницу от пользователя
  document.querySelector('.loader').classList.remove('hidden');
  
  //ставим таймер
  var seconds = 0;
  const timer = setInterval(()=>
  {
    seconds++;
  }, 10);

  //получаем адрес и тип сортировки
  root = document.getElementById('root').value;
  sortType = document.getElementById('sortType')

  //создаем запрос
  let xhr = new XMLHttpRequest();
  let url = new URL(`http://${window.location.host}/dir`);
  
  //устанавливаем в запрос параметры
  url.searchParams.set('root', root);
  if(document.getElementById('sortType').checked == true){
    url.searchParams.set('sortType', 'asc');
  }
  else url.searchParams.set('sortType', 'desc');

  //отправляем запрос
  xhr.open('GET', url);
  xhr.send();

  //если получили ответ
  xhr.onload = function() {
    renderDir(xhr);

    //останавливаем таймер
    const divTimer = document.getElementById('timer');
    clearInterval(timer)
    divTimer.innerHTML=`Время выполнения: ${seconds/100} секунд(ы)`;
  };

  //если получили ошибку
  xhr.onerror = function() { 
    alert('[Ошибка соединения]');
  };
};

//renderDir 
function renderDir(xhr){
  //делаем спинер невидимым
  document.querySelector('.loader').classList.add('hidden');

  //получаем и чистим контейнер со списком файлов, папок
  const divUnswers = document.getElementById('unswers');
  divUnswers.innerHTML = "";

  //обрабатываем ответ сервера в виде json-файла
  unmarshFiles = JSON.parse(xhr.response);
  if(unmarshFiles == null){
    divUnswers.innerText = "Папка пуста";
    return;
  }

  //создаём маркерованный список
  let ul = document.createElement('ul');
  ul.className = "filesUl";
  ul.id = "files";

  //присваиваем контейнеру список
  divUnswers.appendChild(ul);

  //обрабатываем полученные данные
  for (let i = 0; i < unmarshFiles.length; i++){

    //добавляем в список элементы li и в зависимости от типа элемента списка
    //(папка, файл) ставим свою марку
    var li = document.createElement("li");
    li.id = `filesLi${i}`;
    if(unmarshFiles[i].Type == "file") li.className = "lis fileLi"
    else li.className = "lis folderLi"

    //элементу списка li присваем строку, которой обрезаем адрес корневой папки,
    //адрес корневой папки выводится наверху страницы
    li.innerHTML = `${unmarshFiles[i].Name.slice(root.length + 1)}&nbsp;:&nbsp${unmarshFiles[i].Size} байт(ов)`;

    //ставим обработчки нажатия на список ul
    ul.onclick = (event) => {
      //если нажали не на элемент списка, то ничего не делаем
      if(event.target == '[object HTMLUListElement]') return;

      //получаем все элементы списка li
      var dots = document.getElementsByClassName('lis');

      //получаем json-объект по индексу, найденному по индексу элемента li, на который нажали
      clickedFileLi = unmarshFiles[Array.from(dots).indexOf(event.target)]

      //нажатие на файл не обрабатываем
      if(clickedFileLi.Type == 'file') return

      //меняем строку ввода, если папка, и вызываем getDir()
      document.getElementById('root').value = clickedFileLi.Name;
      getDir()
    }
    //присваиваем списку ul элемент li 
    ul.append(li);
  }
}
