<?php
$host = 'localhost'; // адрес сервера
$db_name = 'my_database'; // имя базы данных
$user = 'just_user'; // имя пользователя
$password = '254812User$'; // пароль

//получаем данные
$body = file_get_contents('php://input');
$jsondata = json_decode($body, true);
$root = $jsondata['root'];
$time = $jsondata['elapsed_time'];

//отправляем ответ
http_response_code(201);
$responseData = array(
    'time' => "$time",
    'root' => "$root"
);
echo json_encode(($responseData));

// создание подключения к базе   
   $connection = mysqli_connect($host, $user, $password, $db_name);

// получаем дату
    $today = date("Y-m-d H:i:s");

// текст SQL запроса, который будет передан базе
   $query = "insert into statistics (root, date, elapsed_time) values ('$root', '$today', '$time')";

// выполняем запрос к базе данных
   $result = mysqli_query($connection, $query);

// закрываем соединение с базой
   mysqli_close($connection);
?>
