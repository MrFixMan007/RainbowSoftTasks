import * as rendering from "./script";
export default class DirWorker{
    private defaultRoot : string;

    private loaderId : string;
    private sortTypeCheckboxId : string;
    private timerId : string;
    private root : HTMLLabelElement | null;

    constructor(loaderId : string, rootId : string, defaultRoot : string, sortTypeCheckboxId : string, timerId : string){
        this.defaultRoot = defaultRoot;

        this.loaderId = loaderId;
        this.sortTypeCheckboxId = sortTypeCheckboxId;
        this.timerId = timerId;
        this.root = document.getElementById(rootId) as HTMLLabelElement
    }
//getDir отправляет запрос серверу, обрабатывает ответы, замеряет время выполнения
getDir() { 
    //ставим спинер, скрывая страницу от пользователя
    const loader : HTMLElement = <HTMLElement> document.getElementById(this.loaderId);
    loader.classList.remove('hidden')

    //ставим таймер
    let seconds : number = 0;
    const timer : ReturnType<typeof setInterval> = setInterval(()=>
    {
        seconds++;
    }, 10);

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
        const unmarshFiles : JSONFile[] = JSON.parse(xhr.response);

        //вызываем рендер
        rendering.render.render(unmarshFiles);

        //останавливаем таймер
        const divTimer : HTMLInputElement = <HTMLInputElement> document.getElementById(this.timerId);
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