import * as dirworker from "./dirworker";

const defaultRoot : string = '/home' //значение root по умолчанию 
//получаем адрес корневой папки root и задаём значение по умолчанию
var root : HTMLElement = <HTMLElement> document.getElementById('root')
if(root) root.innerHTML = defaultRoot

//вешаем обработчики
let backButton : HTMLElement = <HTMLInputElement> document.getElementById('backButton')
let sortType : HTMLElement = <HTMLInputElement> document.getElementById('sortType')

//пока данные грузятся, страницу перекрывает спинер
var spinnerLoadDir : HTMLElement = document.createElement('div')
spinnerLoadDir.className = "loader"
spinnerLoadDir.id = "loader"
document.body.append(spinnerLoadDir)

let dirWorker : dirworker.DirWorker = new dirworker.DirWorker(spinnerLoadDir.id, 'unswers', root.id, defaultRoot, 'sortType', 'timer')
export{dirWorker}

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