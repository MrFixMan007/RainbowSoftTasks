import * as dirWorker from "./script";
//renderDir отрисовывает новый список с файлами и папками
export default class RenderDir{
  readonly arrayOfSizeDescription : string[] = [
    "байт",
    "Кбайт",
    "Мбайт",
    "Гбайт",
    "Ебайт"
];
    private loaderId : string;
    private divUnswersId : string;
    private rootId : string;
    private sortByFoldersIdCheckbox : string | null;
    private sortbyFilesIdCheckbox : string | null;
    private unmarshFiles : JSONFile[] | null = null;

    constructor(loaderId : string, divUnswersId : string, rootId : string, sortByFoldersIdCheckbox : string | null = null, sortbyFilesIdCheckbox : string  | null = null) {
      this.loaderId = loaderId
      this.divUnswersId = divUnswersId
      this.rootId = rootId
      this.sortByFoldersIdCheckbox = sortByFoldersIdCheckbox
      this.sortbyFilesIdCheckbox = sortbyFilesIdCheckbox
    }
    private compareByFiles(a : JSONFile, b : JSONFile) {
      if (a.Type == "file") {
        return -1;
      }
      if (a.Type == "folder") {
        return 1;
      }
      return 0;
    }
    private compareByFolders(a : JSONFile, b : JSONFile) {
      if (a.Type == "folder") {
        return -1;
      }
      if (a.Type == "file") {
        return 1;
      }
      return 0;
    }
    
    sortByType(){
      if (this.sortByFoldersIdCheckbox == null || this.sortbyFilesIdCheckbox == null) return;
      let type : string = "";
      const sortByFoldersCheckbox : HTMLInputElement = document.getElementById(this.sortByFoldersIdCheckbox) as HTMLInputElement
      const sortbyFilesCheckbox : HTMLInputElement = document.getElementById(this.sortbyFilesIdCheckbox) as HTMLInputElement
      if (sortByFoldersCheckbox.checked == true) {
        sortbyFilesCheckbox.checked = false;
        alert("folders")
        type = "folder"
      } else if (sortbyFilesCheckbox.checked == true) {
        sortByFoldersCheckbox.checked = false;
        alert("files")
        type = "files"
      }
      switch(type){
        case "":
          return;
        case "folder":
          this.unmarshFiles?.sort(this.compareByFolders)
          if(this.unmarshFiles) this.render(this.unmarshFiles)
        case "file":
          this.unmarshFiles?.sort(this.compareByFiles)
          if(this.unmarshFiles) this.render(this.unmarshFiles)
      }
    }
    render(unmarshFiles : JSONFile[]){
      this.unmarshFiles = unmarshFiles;
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
          let size : number = +unmarshFiles[i].Size
          let countOfDivisionOfSize = 0;
          while(size >= 1024 || countOfDivisionOfSize > this.arrayOfSizeDescription.length){
              size /= 1024
              countOfDivisionOfSize++
          }
          li.innerHTML = `${unmarshFiles[i].Name.slice(rootString.length + 1)}&nbsp;:&nbsp${size.toFixed(1)} ${this.arrayOfSizeDescription[countOfDivisionOfSize]}`;
      
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