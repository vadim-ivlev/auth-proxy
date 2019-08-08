# auth-proxy -  Пользовательская авторизация

__Приложение__

http://auth-proxy.rg.ru/

__Схема__

http://auth-proxy.rg.ru/testapp

__Тестовое GUI приложение__

http://auth-proxy.rg.ru/testapp


------------------------------


Локальные адреса


- Тестовая страница API: <http://localhost:4000/> `GET`.

- Конечная точка GraphQL <http://localhost:4000/schema> `POST`.



Запуск Postgres
-----------------------   

Если не используется SQLite перед запуском приложения нужно запустить Postgres

    docker-compose up -d    



Запуск приложения (для разработчиков)
-----------------

    go run main.go -serve 4000 -env=dev



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


Бенчмарки соединений с БД с пулом и без
--------------------------------------

    go test -run=Bench -benchmem -benchtime=1s -bench=. ./model/db

результаты

    Benchmark_local_DB-4           500     3210246 ns/op    18911 B/op     303 allocs/op
    Benchmark_local_DB_pool-4    10000      221768 ns/op     1181 B/op      28 allocs/op
    Benchmark_remote_DB-4          200     7894973 ns/op    22676 B/op     312 allocs/op
    Benchmark_remote_DB_pool-4   10000      210710 ns/op     1181 B/op      28 allocs/op
    Benchmark_SQLite-4           10000      114440 ns/op     4328 B/op      89 allocs/op
    Benchmark_SQLite_pool-4      10000      209662 ns/op     1511 B/op      29 allocs/op

оптимизация `getKeysAndValues()`

    Benchmark_getKeysAndValues-4              500000              2947 ns/op             848 B/op         28 allocs/op
    Benchmark_getKeysAndValues1-4            1000000               853 ns/op             400 B/op         11 allocs/op



Сборка для фронтэнд разработчиков
----------------------------------------------


    sh/build-frontend-container.sh

или 

    sh/build-frontend-container.sh


# REWRITE

---------------------


О программе
=====================

Данные
-------


<img src="templates/db.png">

Таблицы БД восстанавливаются и наполняются тестовыми данными при каждом запуске приложения.

Для обеспечения ссылочной целостности на таблицы наложены ограничения внешних ключей с каскадным удалением из подчиненных таблиц. На ключи построены индексы.




Файлы и директории
-------------------



    configs/

Содержит настроечные файлы соединений с Postgres и SQLite. 



    controller/


содержит функции GraphQL и REST API.




    middleware/

Каждый запрос к программе обрабатывается двумя функциями middleware до того как будет обработан основным контроллером.

```
Жизненный цикл запроса

(req) --> HeadersMiddleware --> RedisMiddleware --> router --> controller --> (resp)
```
`HeadersMiddleware` добавляет CORS заголовки.  `CheckUser` проверяет пользователя перед перенаправлением запросов к проксируемым приложениям.




    model/
        db/         - работа с базой данных
        auth/       - работа с пользователями
        mail/       - почта
        session/    - работа с сессиями пользователей



    router/

Сопоставляет маршруты функциям-контроллерам, присоединяет middleware, и запускает сервер. 

    migrations/
        

SQL скрипты для порождения объектов базы данных. Миграции исполняются  при каждом запуске программы, поэтому программа будет корректно работать даже при изначально пустой базе данных.



**Второстепенные файлы**


    etc/
        .pgpass

Файл используется контейнером db Postgres, чтобы не вводить пароли при дампе и восстановлении базы данных.


    templates/


Шаблоны приветственного сообщения приложения и тестовой страницы API <http://localhost:4000/>.





    docker-compose.yml     
    main.go
    README.md                           # Этот файл
    build-frontend-container.sh*        # Скрипт сборки докер контейнера для фронтэнд разработчиков
    build-frontend-container-bare.sh*   # Скрипт сборки докер контейнера для фронтэнд разработчиков
    Dockerfile-frontend                 # Используется в build-frontend-container.sh
    Dockerfile-frontend-bare            # Используется в build-frontend-container.sh
    docker-compose-frontend.yml         # Файл запуска для фронтэнд разработчиков. 
    docker-compose-frontend-bare.yml    # Файл запуска для фронтэнд разработчиков. 
    TODO.md                             # Недоделки








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
- Database: auth




Останов базы данных
    
    docker-compose down



Удаление файлов базы данных после останова docker-compose

    sudo rm -rf  pgdata



Дамп базы данных в файл в директорию `migrations/`.
  
    docker exec -it psql-com pg_dump --file /dumps/auth-dump.sql --host "localhost" --port "5432" --username "root"  --verbose --format=p --create --clean --if-exists --dbname "auth"


Восстановление БД из дампа в `migrations/`.

    docker exec -it psql-com psql -U root -1 -d auth -f /dumps/auth-dump.sql



Дамп схемы БД

    docker exec -it psql-com pg_dump --file /dumps/auth-schema.sql --host "localhost" --port "5432" --username "root" --schema-only  --verbose --format=p --create --clean --if-exists --dbname "auth"


Дамп только данных таблиц.

    docker exec -it psql-com pg_dump --file /dumps/auth-data.sql --host "localhost" --port "5432" --username "root"  --verbose --format=p --dbname "auth" --column-inserts --data-only --table=broadcast --table=post --table=image


Можно добавить  -$(date +"-%Y-%m-%d--%H-%M-%S") к имени файла для приклеивания штампа даты-времени.



Показ структуры таблицы TABLE_NAME

    docker-compose exec db pg_dump -U root -d auth -t TABLE_NAME --schema-only



Командная строка Postgres

	docker-compose exec db psql -U root auth



Командная строка Redis

    docker-compose exec redis redis-cli



## контроль деплоя на works

    ssh -i ~/.ssh/deploy_gitupdater_works_open_ssh gitupdater@212.69.111.246

    sudo docker network ls
    sudo docker exec -it auth-proxy-prod bash
    sudo docker logs -f auth-proxy-prod
    sudo docker logs -f auth-node







