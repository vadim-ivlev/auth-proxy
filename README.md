# Пользовательская авторизация

Приложение для авторизации пользователей.


- Тестовая страница API: <http://localhost:4000/> `GET`.

- Конечная точка GraphQL <http://localhost:4000/graphql> `POST`.


Требования к ПО
--------------

На компьютере разработчика должны быть установлены Docker, Docker-compose, Go.



Запуск Postgres и Redis
-----------------------   

Перед запуском приложения необходимо запустить Postgres и Redis

    docker-compose up -d    



Запуск приложения (для разработчиков)
-----------------

    go run main.go -serve 4000



Для просмотра списка возможных параметров запустите программу без параметров.

    go run main.go



## Миграции

**Важно!** При запуске программы запускаются миграции => вся работа с базой данных должна проходить с помощью [миграций](https://github.com/golang-migrate/migrate). Файлы находятся в директории `migrations/`.

**Create**  

    migrate create -ext sql -seq -digits 2 -dir migrations name


**Up, Down, Version, Goto...**  

    migrate -source=file://migrations/ -database postgres://root:root@localhost:5432/auth?sslmode=disable up 
    migrate -source=file://migrations/ -database postgres://root:root@localhost:5432/auth?sslmode=disable down
    migrate -source=file://migrations/ -database postgres://root:root@localhost:5432/auth?sslmode=disable version
    migrate -source=file://migrations/ -database postgres://root:root@localhost:5432/auth?sslmode=disable goto 2



Тесты
--------

Запуск всех тестов

    go test -v ./...


Функциональные тесты (End to End) проводятся с помощью <https://graphql-test.now.sh/>


# REWRITE


Сборка для фронтэнд разработчиков
----------------------------------------------


    ./build-frontend-container.sh

Скрипт `build-frontend-container.sh` генерирует образ приложения согласно `Dockerfile-frontend` 
и выгружает его в <https://hub.docker.com/>. Файл `docker-compose-frontend.yml` ссылается на этот образ 
и служит для запуска приложения фронтэнд разработчиками на локальных компьютерах. 

Порядок запуска приложения фронтэнд разработчиками описан в файле <readme-frontend.md>.



Деплой
-------

Описание настроек для размещения программы на боевых серверах находится в файле <readme-production.md>.






---------------------


О программе
=====================

Данные
-------

Для хранения данных используется СУБД Postgres v11. 


<img src="docs/tables.png">

Таблицы БД восстанавливаются и наполняются тестовыми данными при каждом запуске приложения.
- broadcast  - трансляции
- post  - посты к трансляциям. Таблица рекурсивно ссылается на саму себя для организации ответов к посту.
- image - изображения к постам

Для обеспечения ссылочной целостности БД на таблицы `post` и `image` наложены ограничения внешних ключей с каскадным удалением из подчиненных таблиц. Первичные ключи `id` автогенерируются. На все ключи построены индексы.


Описание полей таблиц находятся в файле `controller/controller-graphql.go`. Во время исполнения могут быть получены стандартными средствами GraphQL. 




Файлы и директории
-------------------



    configs/

Содержит настроечные файлы соединений с Postgres и Redis. Файл `routes.yaml` описывает маршруты и возможные параметры запросов для тестовой страницы API <http://localhost:7777/>. Файл `routes-front.yaml` - усеченная версия файла `routes.yaml`, содержит маршруты публичного REST API для показа на страницах RG.RU, подключается вместо  `routes.yaml` если при запуске программы задан параметр `-front`.




    controller/

Файл `controller-graphql.go` содержит функции GraphQL API, `controller-rest.go` функции REST API.




    middleware/

Каждый запрос к программе обрабатывается двумя функциями middleware до того как будет обработан основным контроллером.

```
Жизненный цикл запроса

(req) --> HeadersMiddleware --> RedisMiddleware --> router --> controller --> (resp)
```
`HeadersMiddleware` добавляет json и возможно CORS HTTP заголовки к ответу сервера.  `RedisMiddleware` кэширует ответы сервера для публичных REST маршрутов начинающихся с `/api/` и возвращает кэшированные ответы клиенту.




    model/
        db/         - работа с базой данных
        redis/      - кэширование в  Redis
        img/        - масштабирование и сохранение загруженных изображений
        imgserver/  - перемещение изображений по ssh


Приложение не использует ORM. JSON форматирование в части REST API возложено на представления и функции postgres. Таким образом снижается трафик между БД и приложением, ускоряются запросы, снижается нагрузка на  хостирующий сервер, и уменьшается объем кода. 

Для формирования JSON в части GraphQL используется библиотека <github.com/jmoiron/sqlx>, что позволило применить отображения (map) вместо структур и избавиться от описания типов БД в коде Go. Это уменьшает связность  БД и приложения, позволяет менять структуру базы данных без необходимости вносить изменения в код go приложения и уменьшает объем кода.



    router/

Сопоставляет маршруты функциям-контроллерам, присоединяет middleware, и запускает сервер. Сами маршруты с именами соответствующих функций вынесены в настроечные файлы `configs/routes.yaml` и `routes-front.yaml`.


    migrations/
        

SQL скрипты для порождения объектов базы данных. Миграции исполняются  при каждом запуске программы, поэтому программа будет корректно работать даже при изначально пустой базе данных.



**Второстепенные файлы**


    etc/
        .pgpass

Файл используется контейнером db Postgres, чтобы не вводить пароли при дампе и восстановлении базы данных.


    templates/


Шаблоны приветственного сообщения приложения и тестовой страницы API <http://localhost:7777/>.


    docs/

Файлы для документации и проч.


    pgdata/

Директория где postgres хранит файлы базы данных. Может быть удалена. Восстанавливается при каждом новом запуске приложения.


    uploads_temp/

Временная директория для хранения загруженных изображений.

    uploads_temp/

Директория для хранения загруженных изображений. Разделяется с 


    docker-compose.yml     
    main.go
    readme-frontend.md           # Для фронтэнд разработчиков
    readme-production.md         # Для админов
    README.md                    # Этот файл
    build-frontend-container.sh* # Скрипт сборки докер контейнера для фронтэнд разработчиков
    Dockerfile-frontend          # Используется в build-frontend-container.sh
    docker-compose-frontend.yml  # Файл запуска для фронтэнд разработчиков. 
    TODO.md                      # Недоделки








-------------------------------------------------------

Другие команды
--------------------


Просмотр состояния базы данных


Postgres доступен на localhost:5432.

Если блок adminer раскомментирован в `docker-compose.yml`, то в браузере откройте <http://localhost:8080>. 

Параметры доступа:
- System: PostgreSQL,
- Server: db,
- Username: root,
- Password: root,
- Database: onlinebc




Останов базы данных
    
    docker-compose down



Удаление файлов базы данных после останова docker-compose

    sudo rm -rf  pgdata



Дамп базы данных в файл в директорию `migrations/`.
  
    docker exec -it psql-com pg_dump --file /dumps/onlinebc-dump.sql --host "localhost" --port "5432" --username "root"  --verbose --format=p --create --clean --if-exists --dbname "onlinebc"


Восстановление БД из дампа в `migrations/`.

    docker exec -it psql-com psql -U root -1 -d onlinebc -f /dumps/onlinebc-dump.sql



Дамп схемы БД

    docker exec -it psql-com pg_dump --file /dumps/onlinebc-schema.sql --host "localhost" --port "5432" --username "root" --schema-only  --verbose --format=p --create --clean --if-exists --dbname "onlinebc"


Дамп только данных таблиц.

    docker exec -it psql-com pg_dump --file /dumps/onlinebc-data.sql --host "localhost" --port "5432" --username "root"  --verbose --format=p --dbname "onlinebc" --column-inserts --data-only --table=broadcast --table=post --table=image


Можно добавить  -$(date +"-%Y-%m-%d--%H-%M-%S") к имени файла для приклеивания штампа даты-времени.



Показ структуры таблицы TABLE_NAME

    docker-compose exec db pg_dump -U root -d onlinebc -t TABLE_NAME --schema-only



Командная строка Postgres

	docker-compose exec db psql -U root onlinebc



Командная строка Redis

    docker-compose exec redis redis-cli











