import * as dirWorker from "./script";
//renderDir отрисовывает новый список с файлами и папками
export default class RenderDir{
    private loaderId : string;
    private divUnswersId : string;
    private rootId : string;

    constructor(loaderId : string, divUnswersId : string, rootId : string) {
      this.loaderId = loaderId
      this.divUnswersId = divUnswersId
      this.rootId = rootId
    }
    render(unmarshFiles : JSONFile[]){
        //убираем спинер (делаем невидимым)
        const loader : HTMLElement = <HTMLElement> document.getElementById(this.loaderId);
        loader.classList.add('hidden')

        const root : HTMLElement = <HTMLElement> document.getElementById(this.rootId);
      
        //получаем и чистим контейнер со списком файлов, папок
        const divUnswers : HTMLElement = <HTMLElement> document.getElementById(this.divUnswersId);
        divUnswers.innerHTML = "";
      

        if(unmarshFiles == null){
          divUnswers.innerText = "Папка пуста";
          return;
        }
      
        //создаём маркированный список
        const ul : HTMLElement = document.createElement('ul');
        ul.className = "filesUl";
        ul.id = "files";
      
        //присваиваем контейнеру список
        divUnswers.appendChild(ul);
      
        //обрабатываем полученные данные
        for (let i = 0; i < unmarshFiles.length; i++){
      
          //добавляем в список элементы li и в зависимости от типа элемента списка
          //(папка, файл) ставим свою марку
          const li : HTMLElement = document.createElement("li");
          li.id = `filesLi${i}`;
          if(unmarshFiles[i].Type == "file") li.className = "lis fileLi"
          else li.className = "lis folderLi"
      
          //элементу списка li присваем строку, которой обрезаем адрес корневой папки,
          //адрес корневой папки выводится наверху страницы
          const rootString : string = String(root?.textContent)
          li.innerHTML = `${unmarshFiles[i].Name.slice(rootString.length + 1)}&nbsp;:&nbsp${unmarshFiles[i].Size} байт(ов)`;
      
          //ставим обработчки нажатия на список ul
          ul.onclick = (event) => {
            const eventTarget : HTMLElement = <HTMLElement> event.target 
            //если нажали не на элемент списка, то ничего не делаем
            const eventTargetString : string = String(event.target) 
            if(eventTargetString == '[object HTMLUListElement]') return;
      
            //получаем все элементы списка li
            const dots : HTMLCollection = document.getElementsByClassName('lis');
      
            //получаем json-объект по индексу, найденному по индексу элемента li, на который нажали
            const clickedFileLi : JSONFile = unmarshFiles[Array.from(dots).indexOf(eventTarget)]
      
            //нажатие на файл не обрабатываем
            if(clickedFileLi.Type == 'file') return
      
            //меняем строку ввода, если папка, и вызываем getDir()
            root.innerHTML = clickedFileLi.Name;
            dirWorker.dirWorker.getDir()
          }
          //присваиваем списку ul элемент li 
          ul.append(li);
        }
      }  
}