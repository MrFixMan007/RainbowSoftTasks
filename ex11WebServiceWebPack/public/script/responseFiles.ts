//класс ResponseFiles для обработки ответа с сервера
class ResponseFiles {
    Time : string;
    Files : JSONFile[] = [];
    constructor(Time : string, Type : string, Name : string, Size : string) {
      this.Time = Time;
      this.Files[this.Files.length] = new JSONFile(Type, Name, Size)
    }
  }