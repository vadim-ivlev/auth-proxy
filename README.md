# auth-proxy -  Пользовательская авторизация



## Мотивация

Web  приложения  нуждаются в **контроле доступа к данным**. 

Для примера представим **редакцию онлайн газеты**. В редакции есть несколько отделов: новости, спорт, политика, технологии, литературная жизнь т.п.. 

Пусть **задача** состоит в создании **web приложения для литературного отдела**.  (API  и HTML). 





![alt_text](templates/images/Sketches-1.png "users and roles")


**Пользователи** приложения делятся на несколько  категорий, с каждой из которых сопоставлена определенная **роль**. Авторы пишут статьи, редакторы правят их, шефы решают, какие материалы будут опубликованы, читатели -  самая многочисленная  роль, не пишут статей, но могут оставлять комментарии к ним. Один и тот же человек может выполнять несколько ролей.

Учет пользователей и их ролей, аутентификация,  авторизация (наделение пользователей ролями, и тем самым делегация им определенными прав) составляет значительную часть приложения.

В настоящей редакции могут быть десятки подобных приложений. Легко представить, как **возрастает сложность** информационной системы редакции с ростом числа приложений, ролей и пользователей. 


## Борьба со сложностью

Идея состоит в том, чтобы вынести функции защиты в отдельное приложение **auth-proxy**, обслуживающее частные приложения системы.



![alt_text](templates/images/Sketches-0.png "image_tooltip")


Таким образом мы избавляем приложения от необходимости заботиться о задачах контроля доступа, уменьшаем код приложений, объемы данных с которыми они работают, не говоря уже о регламентах администрирования пользователей и прочих бизнес правилах.

Клиентами **auth-proxy** могут выступать не только пользователи но и другие приложения.


## Как работает auth-proxy

Жизненный цикл запросов клиента показан на рисунке




![alt_text](templates/images/Sketches-3.png "image_tooltip")


Сначала клиент посылает **auth-proxy** имя и пароль и в случае успешной аутентификации получает в ответ **токен**  сгенерированный **auth-proxy** специально для данного клиента.  Токен служит временным удостоверением для клиента  на текущий сеанс связи. Клиент обязан присоединять его к заголовку каждого последующего запроса. Получая валидный токен сервер убеждается в том, что запрос пришел от аутентифицированного пользователя. 

В целом схема похожа систему контроля допуска на вечеринку. При входе приглашенный называет фамилию и имя, и если он есть в списке ему выдается браслет или другой опознавательный знак, который он должен предъявлять всякий раз, когда решит временно отлучиться.

При следующих запросах клиента **auth-proxy** присоединяет к заголовку запроса список ролей клиента в целевом приложении и пересылает запрос приложению.

После получения ответа приложения auth-proxy пересылает ответ обратно клиенту.


## Разделение ответственности

**Auth-proxy  отвечает, за:**



1. Аутентификацию (идентификацию) пользователей - удостоверение в том, что пользователь именно тот за кого себя выдает. Производится с помощью имени/пароля и web - токенов.
2. Авторизацию (делегирование прав) - определение списка ролей пользователя для каждого конкретного приложения.
3. Проксирование - пересылка запросов пользователя конечному приложению и обратно. 
4. Создание, редактирование, удаление записей о пользователях, ролях и приложениях.

С академической точки зрения,  auth-proxy нужно разделить на 4 отдельных примитивных приложения, каждое из которых выполняет ровно одну задачу. Но из практических соображений четыре приложения были объединены в одно.

**Конечное приложение отвечает за:**



1. За обработку запросов пользователя с учетом его ролей, перечисленных в заголовке запроса


## Модель данных

Ядром авторизации служит таблица **app_user_role**, с полями usernamе, appname, rolename для хранения идентификатора пользователя, приложения и его роли в этом приложении. Пользователь может выполнять несколько ролей в конкретном приложении. В этом случае под каждую роль в таблице заводится отдельная запись.



![alt_text](templates/images/Sketches-4.png "db schema")


Кроме главной таблицы **app_user_role** должны быть три справочные таблицы  user, app и role, для хранения дополнительных сведений о пользователях, приложениях и ролях.

“Должны быть” потому, что ответственность за  смысл (семантику)  набора конкретных привилегий, связанных с ролью, лежит на конкретном приложении, и auth-proxy не обязано об этом ничего знать. По этой причине справочная таблица role отсутствует. Для auth-proxy роль - это просто ничего не значащий идентификатор.



![alt_text](templates/images/Sketches-5.png "image_tooltip")


Таблица таблица `user` используется для аутентификации пользователя , где помимо идентификатора хранится хэш пароля. Сам пароль не сохраняется из соображений секретности.

Таблица `app` используется для проксирования запросов к конечным приложениям и помимо прочего содержит IP адрес конечного приложения, предположительно недоступный с компьютера клиента.


## Токены

Токены выдаются успешно аутентифицированному клиенту и работают по схеме web-токенов (jwt). Токен состоит из тела, подписанного цифровой подписью, сгенерированной Auth-proxy c с помощью секретного ключа. Получив от клиента токен сервер проверяет соответствие подписи телу токена с помощью того же ключа. Таким образом происходит проверка, что токен действительно был сгенерирован Auth-proxy. Если клиентом auth-proxy является веб- приложение, токен сохраняется на компьютере клиента как куки браузера.


## Администрирование пользователей и приложений

**authadmin** - единственная роль определенная в  приложении auth-proxy.

Пользователь с такой ролью имеет максимальные права и может добавлять приложения, назначать роли пользователям и т.п., и может полностью переконфигурировать систему. 

Остальные пользователи могут только



*   просматривать список конечных приложений
*   делать запросы к конечным приложениям 
*   и изменять свои личные данные.


### Замечания

Рекомендуется иметь хотя бы двух пользователей с ролью **authadmin** на случай если один, по ошибке заблокирует, лишит себя этой роли или даже удалит о себе запись.

Вопрос о привилегиях authadmin остается открытым для обсуждения. Возможно имеет смысл запретить ему удаление пользователей, или запретить модификацию записей других пользователей с ролью authadmin, или запретить самоудаление/самоблокирование/саморазроливание.












----------------------------------------------------------

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







