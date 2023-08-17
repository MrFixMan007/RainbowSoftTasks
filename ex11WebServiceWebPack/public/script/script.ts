import DirWorker from "./dirworker";
import RenderDir from "./render";

const defaultRoot : string = '/home' //значение root по умолчанию 
//получаем адрес корневой папки root и задаём значение по умолчанию
let root : HTMLElement = <HTMLElement> document.getElementById('root')
if(root) root.innerHTML = defaultRoot

//вешаем обработчики
let backButton : HTMLElement = <HTMLInputElement> document.getElementById('backButton')
let sortType : HTMLElement = <HTMLInputElement> document.getElementById('sortType')

//пока данные грузятся, страницу перекрывает спинер
let spinnerLoadDir : HTMLElement = document.createElement('div')
spinnerLoadDir.className = "loader"
spinnerLoadDir.id = "loader"
document.body.append(spinnerLoadDir)

let dirWorker : DirWorker = new DirWorker(spinnerLoadDir.id, root.id, defaultRoot, 'sortType', 'timer')
let render : RenderDir = new RenderDir(spinnerLoadDir.id, 'unswers', root.id);

export{dirWorker, render}

if(backButton) {
  backButton.addEventListener('click', getBackDir)
}
if(sortType) {
  sortType.addEventListener('click', getDir)
}

function getBackDir(){
  dirWorker.getBackDir()
}
function getDir(){
  dirWorker.getDir()
}
dirWorker.getDir()