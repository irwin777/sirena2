# sirena2

## Демон отслеживания тревоги по областям с портала [https://alerts.in.ua/](https://alerts.in.ua/).

Облась настраивается по параметру -obl, список областей можно взять из официального апи [API](https://devs.alerts.in.ua/).
Так же необходимо на портале получить уникальный API ключ.
Пример запуска `go run cmd/sirena2/main.go -key "api key" -oblast 31 -trevoga Sub.mp3 -vidbiy Sub.mp3`
Так же есть параметр -test. С помощью него можно проиграть файл тревоги. 

Для работы нужна библиотека alsa-lib.