//класс JSONValue для обработки ответа с сервера
class JSONValue {
  Type : string;
  Name : string;
  Size : string;
  constructor(Type : string, Name : string, Size : string) {
    this.Type = Type;
    this.Name = Name;
    this.Size = Size;
  }
}
const defaultRoot : string = '/' //значение root по умолчанию 

//вешаем обработчики
var backButton : HTMLElement = <HTMLInputElement> document.getElementById('backButton')
var sortType : HTMLElement = <HTMLInputElement> document.getElementById('sortType')
if(backButton) {
  backButton.addEventListener('click', getBackDir)
}
if(sortType) {
  sortType.addEventListener('click', getDir)
}

//пока данные грузятся, страницу перекрывает спинер
var spinnerLoadDir : HTMLElement = document.createElement('div')
spinnerLoadDir.className = "loader"
document.body.append(spinnerLoadDir)

//получаем адрес корневой папки root и задаём значение по умолчанию
var root : HTMLElement = <HTMLElement> document.getElementById('root')
if(root) root.innerHTML = defaultRoot
getDir()

//getBackDir чистит строку вывода с конца до первого встречного слэша
//и вызывает getDir()
function getBackDir(){
  var rootStr : String = "/"
  if(root) rootStr = String(root.textContent)

  var lastIndexOfSlesh : number = rootStr.lastIndexOf("/")
  if(root) root.textContent = rootStr.slice(0, lastIndexOfSlesh)
  if(root.textContent == "/" || root.textContent == ""){
    root.innerHTML = defaultRoot
  }
  getDir()
}

//getDir отправляет запрос серверу, обрабатывает ответы, замеряет время выполнения
function getDir() {  

  //ставим спинер, скрывая страницу от пользователя
  if(spinnerLoadDir) spinnerLoadDir.classList.remove('hidden')
  
  //ставим таймер
  var seconds : number = 0;
  const timer : ReturnType<typeof setInterval> = setInterval(()=>
  {
    seconds++;
  }, 10);

  //получаем тип сортировки
  let sortType : HTMLInputElement = <HTMLInputElement> document.getElementById('sortType')

  //создаем запрос
  let xhr : XMLHttpRequest = new XMLHttpRequest();
  let url : URL = new URL(`http://${window.location.host}/dir`);
  
  //устанавливаем в запрос параметры (адрес root)
  let roorStr : string;
  if(root) {
    roorStr = String(root.textContent);
    url.searchParams.set('root', roorStr);
  }

  //устанавливаем в запрос параметры (тип сортировки)
  if(sortType.checked == true){
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
    let divTimer : HTMLInputElement = <HTMLInputElement> document.getElementById('timer');
    clearInterval(timer)
    divTimer.innerHTML=`Время выполнения: ${seconds/100} секунд(ы)`;
  };

  //если получили ошибку
  xhr.onerror = function() { 
    alert('[Ошибка соединения]');
  };
};

//renderDir отрисовывает новый список с файлами и папками
function renderDir(xhr : XMLHttpRequest){
  //убираем спинер (делаем невидимым)
  spinnerLoadDir.classList.add('hidden')

  //получаем и чистим контейнер со списком файлов, папок
  const divUnswers : HTMLInputElement = <HTMLInputElement> document.getElementById('unswers');
  divUnswers.innerHTML = "";

  //обрабатываем ответ сервера в виде json-файла
  var unmarshFiles : JSONValue[] = JSON.parse(xhr.response);
  if(unmarshFiles == null){
    divUnswers.innerText = "Папка пуста";
    return;
  }

  //создаём маркированный список
  let ul : HTMLElement = document.createElement('ul');
  ul.className = "filesUl";
  ul.id = "files";

  //присваиваем контейнеру список
  divUnswers.appendChild(ul);

  //обрабатываем полученные данные
  for (let i = 0; i < unmarshFiles.length; i++){

    //добавляем в список элементы li и в зависимости от типа элемента списка
    //(папка, файл) ставим свою марку
    var li : HTMLElement = document.createElement("li");
    li.id = `filesLi${i}`;
    if(unmarshFiles[i].Type == "file") li.className = "lis fileLi"
    else li.className = "lis folderLi"

    //элементу списка li присваем строку, которой обрезаем адрес корневой папки,
    //адрес корневой папки выводится наверху страницы
    let rootString : string = String(root?.textContent)
    li.innerHTML = `${unmarshFiles[i].Name.slice(rootString.length + 1)}&nbsp;:&nbsp${unmarshFiles[i].Size} байт(ов)`;

    //ставим обработчки нажатия на список ul
    ul.onclick = (event) => {
      let eventTarget : HTMLElement = <HTMLElement> event.target 
      //если нажали не на элемент списка, то ничего не делаем
      let eventTargetString : string = String(event.target) 
      if(eventTargetString == '[object HTMLUListElement]') return;

      //получаем все элементы списка li
      var dots : HTMLCollection = document.getElementsByClassName('lis');

      //получаем json-объект по индексу, найденному по индексу элемента li, на который нажали
      let clickedFileLi : JSONValue = unmarshFiles[Array.from(dots).indexOf(eventTarget)]

      //нажатие на файл не обрабатываем
      if(clickedFileLi.Type == 'file') return

      //меняем строку ввода, если папка, и вызываем getDir()
      root.innerHTML = clickedFileLi.Name;
      getDir()
    }
    //присваиваем списку ul элемент li 
    ul.append(li);
  }
}
