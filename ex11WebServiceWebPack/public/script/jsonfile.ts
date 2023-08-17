//класс JSONFile для обработки ответа с сервера
class JSONFile {
    Type : string;
    Name : string;
    Size : string;
    constructor(Type : string, Name : string, Size : string) {
      this.Type = Type;
      this.Name = Name;
      this.Size = Size;
    }
  }