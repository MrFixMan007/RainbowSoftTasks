<!DOCTYPE HTML>
<html>
 <head>
  <meta charset="utf-8">
  <title>Статистика</title>
  <style type="text/css">
    #centerLayer {
     margin: 0 auto; /* Отступ слева и справа */
     text-align: left; /* Выравнивание содержимого слоя по левому краю */
    }
   </style>
 </head>
 <body>
 <a href="http://192.168.81.45:8080/" id="back">Назад</a>
  <table border="1" id="centerLayer">
   <caption>Таблица Статистика</caption>
   <tr>
    <th>Id</th>
    <th>Адрес</th>
    <th>Дата обращения</th>
    <th>Затраченное время</th>
    <th>Размер</th>
   </tr>
   <?php
      $host = 'localhost'; // адрес сервера
      $db_name = 'my_database'; // имя базы данных
      $user = 'just_user'; // имя пользователя
      $password = '254812User$'; // пароль

      // создание подключения к базе   
         $connection = mysqli_connect($host, $user, $password, $db_name);

      // текст SQL запроса, который будет передан базе
         $query = 'SELECT * FROM `statistics`';

      // выполняем запрос к базе данных
         $result = mysqli_query($connection, $query);

      // выводим полученные данные
   while($row = $result->fetch_assoc()){
      echo "<tr>";
      echo "<td>";
      echo  $row['id'];
      echo "</td>";
      echo "<td>";
      echo  $row['root'];
      echo "</td>";
      echo "<td>";
      echo  $row['date'];
      echo "</td>";
      echo "<td>";
      echo  $row['elapsed_time'];
      echo " миллисекунд(ы)";
      echo "</td>";
      echo "<td>";
      echo  $row['size'];
      echo " байт";
      echo "</td>";
      echo "</tr>";
   }

// закрываем соединение с базой
   mysqli_close($connection);
   ?>
  </table>
 </body>
</html>
