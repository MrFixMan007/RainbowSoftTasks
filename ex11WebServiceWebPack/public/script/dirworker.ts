import * as rendering from "./render";
class DirWorker{
    private defaultRoot : string;

    private loaderId : string;
    private rootId : string;
    private divUnswersId :string;
    private sortTypeCheckboxId : string;
    private timerId : string;
    // private root : HTMLLabelElement | null;

    constructor(loaderId : string, divUnswersId : string, rootId : string, defaultRoot : string, sortTypeCheckboxId : string, timerId : string){
        this.defaultRoot = defaultRoot;

        this.loaderId = loaderId;
        this.rootId = rootId;
        this.divUnswersId = divUnswersId;
        this.sortTypeCheckboxId = sortTypeCheckboxId;
        this.timerId = timerId;
        // this.root = document.getElementById(rootId)
    }
//getDir отправляет запрос серверу, обрабатывает ответы, замеряет время выполнения
getDir() {  
    let root : HTMLElement =  <HTMLElement> document.getElementById(this.rootId);
    //ставим спинер, скрывая страницу от пользователя
    let loader : HTMLElement = <HTMLElement> document.getElementById(this.loaderId);
    loader.classList.remove('hidden')

    //ставим таймер
    var seconds : number = 0;
    const timer : ReturnType<typeof setInterval> = setInterval(()=>
    {
        seconds++;
    }, 10);

    //получаем тип сортировки
    let sortType : HTMLInputElement = <HTMLInputElement> document.getElementById(this.sortTypeCheckboxId)

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
    xhr.onload = () => {
        //обрабатываем ответ сервера в виде json-файла
        let unmarshFiles : JSONFile[] = JSON.parse(xhr.response);

        //вызываем рендер
        let render : rendering.RenderDir = new rendering.RenderDir(this.loaderId, this.divUnswersId, this.rootId);
        render.render(unmarshFiles);

        //останавливаем таймер
        let divTimer : HTMLInputElement = <HTMLInputElement> document.getElementById(this.timerId);
        clearInterval(timer)
        divTimer.innerHTML=`Время выполнения: ${seconds/100} секунд(ы)`;
    };

    //если получили ошибку
    xhr.onerror = function() { 
        alert('[Ошибка соединения]');
    };
};

//getBackDir чистит строку вывода с конца до первого встречного слэша
//и вызывает getDir()
getBackDir(){
    let root : HTMLElement =  <HTMLElement> document.getElementById(this.rootId);
    let rootStr : String = "/"
    if(root) rootStr = String(root.textContent)
  
    let lastIndexOfSlesh : number = rootStr.lastIndexOf("/")
    if(root) root.textContent = rootStr.slice(0, lastIndexOfSlesh)
    if(root.textContent == "/" || root.textContent == ""){
        root.innerHTML = this.defaultRoot
    }
    this.getDir()
  }
}
export {DirWorker};