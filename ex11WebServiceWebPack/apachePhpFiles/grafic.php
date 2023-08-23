<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Grafic</title>
    <link rel="stylesheet" href="bower_components/chartist/dist/chartist.min.css">
    <style type="text/css">
      #centerLayer {
     margin: 0 auto; /* Отступ слева и справа */
    }
        /* Используйте этот селектор для переопределения стиля линейного графика */
        .ct-series-a .ct-line {
            /* Цвет линии серии */
            stroke: red;
            /* Толщина линии */
            stroke-width: 5px;
            /* Пунктирная линия по шаблону */
            stroke-dasharray: 10px 20px;
        }

        /* Этот селектор перкрывает стиль точек линейного графика */
        .ct-series-a .ct-point {
            /* Цвет точек */
            stroke: rgb(175, 0, 0);
            /* Размер точек */
            stroke-width: 10px;
            /* Сделать точки квадратами */
            stroke-linecap: square;
        }
    </style>
  </head>
  <?php
  $arrayOfsizeDescription = array(
    0 => "б",
    1 => "Кб",
    2 => "Мб",
    3 => "Гб",
    4 => "Тб"
);

      $host = 'localhost'; // адрес сервера
      $db_name = 'my_database'; // имя базы данных
      $user = 'just_user'; // имя пользователя
      $password = '254812User$'; // пароль

      // создание подключения к базе   
         $connection = mysqli_connect($host, $user, $password, $db_name);

      // текст SQL запроса, который будет передан базе
         $query = 'SELECT * FROM `statistics` order by size';

      // выполняем запрос к базе данных
         $result = mysqli_query($connection, $query);
         $result1 = mysqli_query($connection, $query);
         

// закрываем соединение с базой
   mysqli_close($connection);
   ?>
  <body>
    <h5>График зависимости времени (мс) от размера</h5>
    <!-- Site content goes here !-->
    <script src="bower_components/chartist/dist/chartist.min.js"></script>

    <div class="ct-chart ct-perfect-fourth" id = 'grafic'></div>
    <script type="text/javascript">
        var data = {
            labels: [
              <?php
              while($row = $result->fetch_assoc()){
                echo '\'';
                $size = $row['size'];
                $j = 0;
                while($size >= 1024){
                  $size /= 1024;
                  $j++;
                }
                $roundedSize = round($size);
                echo $roundedSize . " " . $arrayOfsizeDescription[$j];
                echo '\',';
              }
              ?>
              ],
            series: [
                [
                  <?php
              while($row = $result1->fetch_assoc()){
                echo '\'';
                echo $row['elapsed_time'];
                echo '\',';
              }
              ?>
                ]
            ]
        };

        var options = {
          height: '700px'
        };

        new Chartist.Line('.ct-chart', data, options);
    </script>
  </body>
</html>