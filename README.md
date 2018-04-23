# Goque — Микросервис для работы с очередью http/https запросов
**Запуск**
* Для запуска необходим  Docker
* Запущенный контейнер с Redis
* ```git clone https://github.com/legostin/goque.git```
* ```cd goque```
* ```docker --network=<your network> build -f ./Dockerfile .```


**Пример использования в PHP скрипте**
*
```PHP
$data = [
                "url" => "https://google.com",
                "method" => "POST",
                "tag" => "task.test",
                "jsonData" => Json::encode(['foo'=>"bar"])
            ];
            Yii::$app->redis->rpush("tasks",Json::encode($data));
```
*
```PHP
$data = [
                   "url" => "https://google.com",
                   "method" => "GET",
                   "tag" => "task.test",
                   "params" => ["foo"=>"bar"]
               ];
               Yii::$app->redis->rpush("tasks",Json::encode($data));
```
*
```PHP
$data = [
                   "url" => "https://google.com",
                   "method" => "POST",
                   "tag" => "task.test",
                   "params" => ["foo_post"=>"bar"]
               ];
               Yii::$app->redis->rpush("tasks",Json::encode($data));
```
