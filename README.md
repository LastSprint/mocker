# Mocker

Это приложение, который позволяет имитировать работу реального backend'а.

Само по себе приложения представляет из себя обычный ротуер HTTP-запросов. 
Приложение принимает запрос из сети, парсит URL и HTTP-method. Затем проходится по известным мокам. Берет мок с подходящим URL и HTTP-method'ом, читает код ответа и тело ответа, а затем вовзращает ответ клиенту.

В случае, если существует несколько моков с одинаковыми URL и Method, то приложение объединяет их в группу а затем закольцовано итерирует по ним.

После того, как были внесены изменения в моки необходимо, чтобы сервер обновил у себя данные.

Сделать это можно так: `curl 127.0.0.1:port/update_models`

## Моки

Моковые файлы представляют из себя `Json` формата:
```
 {
    "url": string
    "method": string
    "statusCode": int
    "response": object
 }
```

Буквально здесь записано слудеюущее:
На запрос с `URL = url` и `Method = method` вернуть ответ `response` с кодом `statusCode`

### URL

Может быть следующих видов:

1. `path/to/endpoint` - обычный url адрес. Во время получения запроса сервис будет сравнивать строки посимвольно.
2. `path/to/endpoint/{number}` - url с path-паттерном. Мок с таким URL будет реагировать на любой запрос который удовлетворяет этому шаблону.
3. `path/to/endpoint/data?param={value}` - url c parameter-паттером. Мок с таким url отработае на запрос, содержащий заданные параметры. При этом, если одного из параметров не будет в запросе, то он не сматчится с шаблоном.

### Method

Пишите названия HTTP методов в UpperCase. 
**НЕ** `get` **А**  `GET` 

### Status Code

Любое целое число. Желательно из известных [HTTP-кодов](https://ru.wikipedia.org/wiki/Список_кодов_состояния_HTTP)

### Response

Это поле содержит `Json`, который вернется в ответ на запрос. 

## Конфигурация

Приложение имеет конфиг в `Json` следующего формата:
```
{
    "mocksRootDir": string,
    "port": int,
    "logPath": string
}
```

### Mocks Root Dir

Это путь до папки, которая является корневой для файлов с моками. 

### Port

TCP-порт, который будет слушать приложение

### Log Path

Путь до файла с логами. Туда приложение будет писать логи работы.


