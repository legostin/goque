package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"net/url"
)

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	if err !=nil {
		fmt.Println("redis is not ready")
		time.Sleep(time.Second*1)
		NewRedisClient()
	}
	return client

}

func doRequest(req *http.Request, client *redis.Client, requestData string) {
	// Выполняем запрос
	client_http := &http.Client{}
	resp, err := client_http.Do(req)
	if err != nil {
		client.LPush("tasks",requestData)
		main() // Че то не то,  перезапускаем все
	}
	defer resp.Body.Close() //Не забывае закрывать соединение, мы не хотим утечек памяти

	// Удаленный сервер ответил некорректно, ставим данный запрос в начало очереди
	if resp.Status!="200 OK" {
		client.LPush("tasks",requestData)
	}
}

func main() {
	client := NewRedisClient()
	// Поскольку список хэдеров заранее не определен, то определяем интерфейс
	type Headers map[string]interface{}
	// То же самое проделываем и с Параметрами для Гет и Пост запросов
	type Params map[string]interface{}
	// Определяем новый тип для входных данных
	type MyRequest struct {
		Url string // Адрес для запроса
		Method string // Метод Гет или Пост
		Headers Headers // Хэдеры
		JsonData string // Для запросов с content-type:Application/Json
		Params Params // Для остальных запросов
	}
	for {
		fmt.Println("start")
		requestData, err := client.RPop("tasks").Result()
		if err == redis.Nil {
			// Очередь запросов Пуста, ждем секунду, чтобы не перезагружать редис запросами
			time.Sleep(time.Second)
		} else if err != nil {
			// Что-то случилось Редисом, ждем 10 секунд и пытаемся переподключиться
			time.Sleep(time.Second*10)
			main()
		} else {
			var requestStruct MyRequest
			json.Unmarshal([]byte(requestData), &requestStruct) //Маппим строку из редиса в созданную ранее структуру
			fmt.Println(requestStruct.Url)
			if requestStruct.Method=="POST" && requestStruct.Params==nil {
				//Создаем новый запрос. Изначально в body ставим nil
				req, err := http.NewRequest("POST", requestStruct.Url, nil)

				// Если есть данные в джейсоне, то задаем в качестве body эту строку
				if requestStruct.JsonData != "" {
					req.Body = ioutil.NopCloser(bytes.NewBufferString(requestStruct.JsonData))
				}

				//Заполняем хэдеры
				if requestStruct.Headers!=nil {
					for key, value := range requestStruct.Headers {
						fmt.Println(key, value.(string))
						req.Header.Set(key, value.(string))
					}
				}

				if (err!=nil){
					main()
				}
				doRequest(req,client,requestData)

			} else if requestStruct.Method=="POST" && requestStruct.Params!=nil {
				//Создаем новый запрос. Изначально в body ставим nil
				//Добавляем обычные Пост параметры
				form := url.Values{}

				for key, value := range requestStruct.Params {
						form.Add(key, value.(string))
				}
				body := bytes.NewBufferString(form.Encode())

				req, err := http.NewRequest("POST",requestStruct.Url, body)
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

				//Заполняем хэдеры
				if requestStruct.Headers!=nil {
					for key, value := range requestStruct.Headers {
						fmt.Println(key, value.(string))
						req.Header.Set(key, value.(string))
					}
				}
				// Выполняем запрос запрос
				if (err!=nil){
					main()
				}
				doRequest(req,client,requestData)

			} else {

				//Все то же самое и для Гет запроса
				req, err := http.NewRequest("GET", requestStruct.Url, nil)
				if requestStruct.Headers!=nil {
					for key, value := range requestStruct.Headers {
						fmt.Println(key, value.(string))
						req.Header.Set(key, value.(string))
					}
				}
				if requestStruct.Params!=nil {
					requestStruct.Url +="?foo=bar" // небольшой хак, чтобы не париться со строкой
					for key, value := range requestStruct.Params {
						requestStruct.Url += "&"+key+"="+value.(string)
					}
				}
				if (err!=nil){
					main()
				}
				doRequest(req,client,requestData)


			}
		}
	}

}
