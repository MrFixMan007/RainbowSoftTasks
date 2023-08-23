import DirWorker from "./dirworker"
import RenderDir from "./render"

var defaultRoot : string = '/home' //значение root по умолчанию 
//получаем адрес корневой папки root и задаём значение по умолчанию
const root : HTMLElement = <HTMLElement> document.getElementById('root')
if(root) root.innerHTML = defaultRoot

//вешаем обработчики
const backButton : HTMLElement = <HTMLInputElement> document.getElementById('backButton')
const sortType : HTMLElement = <HTMLInputElement> document.getElementById('sortType')
const aStatistic : HTMLLinkElement = <HTMLLinkElement> document.getElementById('aStatistic')

//пока данные грузятся, страницу перекрывает спинер
const spinnerLoadDir : HTMLElement = document.createElement('div')
spinnerLoadDir.className = "loader"
spinnerLoadDir.id = "loader"
document.body.append(spinnerLoadDir)

const dirWorker : DirWorker = new DirWorker(spinnerLoadDir.id, root.id, defaultRoot, 'sortType')
const render : RenderDir = new RenderDir(spinnerLoadDir.id, 'unswers', root.id, 'timer')

export{dirWorker, render}

if(backButton) {
  backButton.addEventListener('click', getBackDir)
}
if(sortType) {
  sortType.addEventListener('click', getDir)
}
if(aStatistic){
  aStatistic.href = `http://${window.location.hostname}:80/getStats.php`
}

function getBackDir(){
  dirWorker.getBackDir()
}
function getDir(){
  dirWorker.getDir()
}
dirWorker.getDir()