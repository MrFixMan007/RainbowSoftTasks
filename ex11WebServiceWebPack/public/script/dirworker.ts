import * as rendering from "./script";
export default class DirWorker{
    private defaultRoot : string;

    private loaderId : string;
    private sortTypeCheckboxId : string;
    private root : HTMLLabelElement | null;

    constructor(loaderId : string, rootId : string, defaultRoot : string, sortTypeCheckboxId : string){
        this.defaultRoot = defaultRoot;

        this.loaderId = loaderId;
        this.sortTypeCheckboxId = sortTypeCheckboxId;
        this.root = document.getElementById(rootId) as HTMLLabelElement
    }
//getDir отправляет запрос серверу, обрабатывает ответы, замеряет время выполнения
getDir() { 
    //ставим спинер, скрывая страницу от пользователя
    const loader : HTMLElement = <HTMLElement> document.getElementById(this.loaderId);
    loader.classList.remove('hidden')

    //получаем тип сортировки
    const sortType : HTMLInputElement = <HTMLInputElement> document.getElementById(this.sortTypeCheckboxId)

    //создаем запрос
    const xhr : XMLHttpRequest = new XMLHttpRequest();
    const url : URL = new URL(`http://${window.location.host}/dir`);

    //устанавливаем в запрос параметры (адрес root)
    let roorStr : string;
    if(this.root) {
        roorStr = String(this.root.textContent);
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
    xhr.onload = () => {
        //обрабатываем ответ сервера в виде json-файла
        const unmarshFiles : ResponseFiles = JSON.parse(xhr.response);
        //вызываем рендер
        rendering.render.render(unmarshFiles);
    };

    //если получили ошибку
    xhr.onerror = function() { 
        alert('[Ошибка соединения]');
    };
};

//getBackDir чистит строку вывода с конца до первого встречного слэша
//и вызывает getDir()
getBackDir(){
    let rootStr : String = "/"
    if(this.root) rootStr = String(this.root.textContent)
  
    const lastIndexOfSlesh : number = rootStr.lastIndexOf("/")
    if(this.root) this.root.textContent = rootStr.slice(0, lastIndexOfSlesh)
    if(this.root?.textContent == "/" || this.root?.textContent == ""){
        this.root.innerHTML = this.defaultRoot
    }
    this.getDir()
  }
}