# auth-proxy -  Пользовательская авторизация

API

- RG production
https://auth-proxy.rg.ru

- RG staging
https://rg-auth-staging.rg.ru


- GL production
https://gl-auth.rg.ru

- GL staging
https://gl-auth-staging.rg.ru


Админки

- RG production
https://auth-proxy.rg.ru/admin/?url=https://auth-proxy.rg.ru

- RG staging
https://rg-auth-staging.rg.ru/admin/?url=https://rg-auth-staging.rg.ru


- GL production
https://gl-auth.rg.ru/admin/?url=https://gl-auth.rg.ru

- GL staging
https://gl-auth-staging.rg.ru/admin/?url=https://gl-auth-staging.rg.ru


<br>

<br>


## Мотивация

**Контроль доступа к данным**.

В **редакции онлайн газеты** есть несколько отделов: новости, спорт, политика, технологии, литературная жизнь т.п..


![alt_text](templates/images/Sketches-1.png "users and roles")

**Пользователи**  делятся на несколько  "категорий", с каждой из которых сопоставлена определенная **роль**. Авторы пишут статьи, редакторы правят их, шефы решают, какие материалы будут опубликованы, читатели не пишут статей, но могут оставлять комментарии. Один и тот же **пользователь может выполнять несколько ролей**.

Учет пользователей, ролей, аутентификация,  авторизация (наделение пользователей ролями, и тем самым делегация им определенными прав) составляет значительную часть приложения.

В редакции могут быть десятки таких приложений. Легко представить, как **возрастает сложность** контроля доступа с ростом числа приложений, ролей и пользователей.


## Борьба со сложностью

Идея -  вынести функции защиты в отдельное приложение **auth-proxy**.


![alt_text](templates/images/Sketches-0.png "image_tooltip")


**Auth- proxy** избавляет приложения от задач контроля доступа, уменьшает код приложений, объем и сложность данных с которыми они работают, регламенты администрирования и другие бизнес правила.

Клиентами **auth-proxy** могут быть не только пользователи, но и  приложения.


## Как работает auth-proxy






![alt_text](templates/images/Sketches-3.png "image_tooltip")


Сначала клиент посылает **auth-proxy** имя и пароль и в случае успешной аутентификации
получает в ответ **токен**  сгенерированный  для данного клиента.

Токен служит временным удостоверением для клиента  на текущий сеанс связи.
Клиент обязан присоединять его к заголовку каждого последующего запроса.
Получая валидный (правильный) токен сервер удостоверяется в том,
что запрос пришел от аутентифицированного пользователя.

Клиент посылает запрос к целевому приложению с присоединенным токеном.

Аuth-proxy присоединяет к заголовку запроса список ролей клиента в целевом приложении и пересылает запрос приложению.

После получения ответа приложения auth-proxy пересылает ответ обратно клиенту.

В целом схема похожа на контроль приглашенных на мероприятие.
При входе гость называет фамилию и имя, и если он есть в списке
ему выдается браслет (баджик/токен), который он должен предъявлять всякий раз,
когда решит вернуться после временной отлучки.




## Разделение ответственности

**Auth-proxy  отвечает, за:**



1. Аутентификацию (идентификацию) пользователей - удостоверение в том, что пользователь именно тот
   за кого себя выдает. Производится с помощью имени/пароля и web - токенов.
2. Авторизацию (делегирование прав) - определение списка ролей пользователя для каждого конкретного приложения.
3. Проксирование — пересылка запросов пользователя конечному приложению и обратно.
4. Создание, редактирование, удаление записей о пользователях, ролях и приложениях.

С академической точки зрения,  auth-proxy нужно бы разделить на 4 отдельных приложения, по одному на задачу,
но из практических соображений они были объединены в одно.

**Конечное приложение отвечает за:**



1. За обработку запросов пользователя с учетом его ролей, перечисленных в заголовке запроса.


## Модель данных

Ядром авторизации служит таблица **app_user_role** (и представление **app_user_role_union**), с полями **usernamе, appname, rolename**
для хранения идентификатора пользователя, имени приложения и роли пользователя в этом приложении.
Пользователь может выполнять несколько ролей в конкретном приложении.
В этом случае под каждую роль в таблице заводится отдельная запись.



![alt_text](templates/images/Sketches-4.png "db schema")


Кроме главной таблицы **app_user_role** должны быть три справочные таблицы  user, app и role, для хранения дополнительных сведений о пользователях, приложениях и ролях. На самом деле их две: **user** и **app**.

Справочная таблица role отсутствует потому, что ответственность за  набор конкретных привилегий, связанных с ролью, лежит на конечном приложении, и auth-proxy не обязано о них знать.  Для auth-proxy роль - это ничего не значащий идентификатор.


![alt_text](templates/images/Sketches-5.png "image_tooltip")

Таблица `user` используется для аутентификации пользователя , где кроме идентификатора хранится хэш пароля. Сам пароль не сохраняется из соображений секретности.

Таблица `app` используется для проксирования запросов к конечным приложениям и помимо прочего содержит IP адрес конечного приложения, возможно недоступный с компьютера клиента.

Полная схема

![alt_text](templates/images/db_schema.png "image_tooltip")




## Токены

Токены выдаются успешно аутентифицированному клиенту и работают по схеме web-токенов (jwt). Токен состоит из тела, подписанного цифровой подписью, сгенерированной Auth-proxy c с помощью секретного ключа. Получив от клиента токен сервер проверяет соответствие подписи телу токена с помощью того же ключа. Таким образом происходит проверка токена на подлинность. Если клиентом auth-proxy является веб - приложение, токен сохраняется на компьютере клиента как куки браузера.


## Администрирование пользователей и приложений

**authadmin** - Пользователь с такой ролью имеет максимальные права и может добавлять приложения, пользователей, назначать роли пользователям и менять параметры работы auth-proxy.

**auditor** - Пользователь с такой ролью может просматривать информацию
о всех пользователях, приложениях и ролях пользователей в приложениях, но не может ничего менять ни в данных ни в параметрах  auth-proxy. Эта роль введена для приложений нуждающихся в информации о пользователях и их ролях
или организующих свою логику распределения пользователй по ролям и группам.


Пользователи без прав **authadmin** и **auditor** могут только


*   просматривать список конечных приложений
*   делать запросы к конечным приложениям
*   и изменять свои личные данные.



Рекомендуется иметь хотя бы двух пользователей с ролью **authadmin**
на случай если один, по ошибке заблокирует, лишит себя этой роли или даже удалит о себе запись.

Вопрос о привилегиях authadmin остается открытым для обсуждения. Возможно имеет смысл
запретить ему удаление пользователей, или запретить модификацию записей других
пользователей с ролью authadmin, или запретить самоудаление/самоблокирование/саморазроливание.


Сброс забытого пароля
----------------------------------------------------------------
```mermaid
sequenceDiagram
    autonumber
    auth-admin ->>+ auth-proxy : reset_password( username, authadmin, authurl )

    auth-proxy ->> DB: ++ user-pashash
    DB ->> auth-proxy: ок
    auth-proxy ->> -auth-admin: ок


    Note over auth-proxy: Посылает  пользователю Email<br>со ссылкой на страницу установки пароля: <br> https://adminurl/set-password.html &pashash &authurl

    Note over auth-admin: Перейдя по ссылке,<br> пользователь делает запрос<br> на установку пароля

    auth-admin -->>+ auth-proxy: set_password(password, pashash)
    auth-proxy -->> DB: -- user-pashash ++ user-password
    DB ->> auth-proxy: ок
    auth-proxy -->>- auth-admin: OK
```
<!-- <img src="templates/images/reset-password.png"> -->


## Oauth2

Приложение допускает аутентификацию и регистрацию пользователей используя сервисы  **facebook, google, yandex, vk** и т.п. Параметры сервисов перечислены в настроечном файле `configs/oauth2.yaml`.

Для аутентификации пользователя через социальную сеть, необходимо, чтобы емайл
пользователя в auth-proxy совпадал с емайлом социальной сети.

Если пользователь не зарегистрирован в auth-proxy, и саморегистрация пользователей
разрешена, после аутентификации пользователя в социальной сети он может зарегистрироваться в auth-proxy. В этом случае сведения из социальной сети будут скопированы в новую запись о пользователе в auth-proxy.

Если пользователь желает зарегистрирваться минуя аутентификацию в социальной сети,
ему придется вводить даннные о себе вручную. В любом случае, перед сохранением записи, пользователь может скорректировать свои данные.

Перечень доступных сервисов Oauth2 с конечными точками  входа и выхода доступен по REST маршруту:

    GET /oauthproviders


И по GraphQL с помощью метода:

    list_oauth_providers()



## Каптча

Чтобы усложнить взлом методом автоматического перебора паролей, в систему встроена **каптча**.

После трех ошибок ввода пароля, сервер начинает требовать
каптчу.

Каптча доступна по URI `/captcha`.

Узнать нужно ли пользователю вводить каптчу можно с помощью метода `is_captcha_required(username)`.

Счетчик неудачных попыток сбрасывается после удачного входа и по истечении часа.

Время сброса счетчика неудачных попыток и максимально допустимое число ошибок
ввода пароля задается через параметры запуска программы:

* -reset_time (в минутах. По умолчанию 60).
* -max_attempts (По умолчанию 5)

## PIN Google Authenticator
Чтобы усложнить взлом методом автоматического перебора паролей, в систему встроен **PIN**.

Узнать нужно ли пользователю вводить PIN можно с помощью метода `is_pin_required(username)`.

## Изменение параметров системы во время исполнения

Параметры приложения перечисленые в файле `congigs/app.env`, считываются в момент
запуска приложения, но некоторые из них допускают изменение во время исполнения с помощью метода `set_params`, право на вызов которого имеет только **authadmin**.

    set_params(
        use_captcha   true  # Нужно ли вводить капчу при входе в систему
        use_pin       false # Нужно ли вводить PIN Google Authenticator при входе в систему
        max_attempts  5     # Максимально допустимое число ошибок ввода пароля
        reset_time    60    # Время сброса счетчика ошибок пароля в минутах
        selfreg       false # Могут ли пользователи регистрироваться самостоятельно
    )

Текущие значения параметров можно получить с помощью метода:

    get_params()


## Контроль состояния системы

Состояние памяти приложения и нагрузку
можно проконтролировать с помощью метода:

    get_stat {
        sys                     # Тотальное количество памяти полученное приложением от OS в мегабайтах
        alloc                   # текущее память аллоцированная из кучи в мегабайтах
        total_alloc             # Кумулятивное количество памяти аллоцированной из кучи в мегабайтах
        requests_per_day        # запросов в день
        requests_per_hour       # запросов в час
        requests_per_minute     # запросов в минуту
        requests_per_second     # запросов в секунду
    }


Те же данные можно получить делая `GET` запрос по маршруту `/stat`.
В последнем случае сам запрос на получение статистики не учитывается при вычислениях количеств запросов за единицу времени.

Маршрут `/metrics` выдает стандартную статистику [Prometheus](https://prometheus.io/)

## Регистронезависимость имени и электронного адреса

При сохранении в базу данных Email и имя пользователя приводятся к нижнему регистру.
При входе в систему регистр букв игнорируется.




## Флаги проксируемых приложений `public`, `rebase`, `sign`



- `public` Установите в "Y" чтобы сделать приложение доступным для
  неавторизованных пользователей.
  В списке приложений пользователь может видеть только те,
  в которых у него определена роль и приложения с `public=Y`.
  Другие приложения недоступны пользователю, как для перечислений в списках поиска,
  так и для проксирования.

- `rebase`
    Поле имеет смысл только для проксирования HTML web приложений.

    Если значение поля равно "Y" приложение пытается "подправить" ссылки, изображения, фреймы в HTML и CSS так, чтобы URI начинающиеся с /, начинались с /apps/appname/. Это делается для того, чтобы браузер не запрашивал ресурсы там, где их нет, то есть в корне auth-proxy.

    Эта операция влечет за собой увеличение потребления памяти и загрузки процессора сервера. Поэтому в приложениях предназначенных для проксирования рекомендуется использовать относительные ссылки и оставлять это поле пустым.
- `sign`
    Если значение поля "Y" запросы к приложению подписываются цифровой подписью в соответствии со спецификацией RFC
    <https://tools.ietf.org/html/draft-cavage-http-signatures-06>.

    Это один из способов контроля источников запросов наряду с ограничениями по IP адресу.
    Позволяет конечному приложению удостовериться, что запрос поступил от auth-proxy.
    Важен для публично доступных в Интернет приложений, которые хотели бы принимать запросы через
    auth-proxy.


    Операция подписи ресурсоемкая. Пользуйтесь разумно. Например не имеет смысла подписывать
    публично доступные приложения или HTML приложения.

    Примеры верификации подписи на golang https://github.com/spacemonkeygo/httpsig, на python https://github.com/spacemonkeygo/httpsig

    Файл `signature/verify.go` содержит пример функции Verify(request) для верификации подписи на стороне конечного приложения.



## Параметры командной строки auth-proxy

- -port default=4400, Запустить приложение на порту 4400
- -config - Конфигурационный файл приложения (default configs/app.env).
- -pgconfig - Конфигурационный файл Postgres (default configs/db.env).



## Конфигурационный файл  app.env

Параметры по умолчанию приложения находятся в в файле  `configs/app.env`.
Файл с параметрами можно указать аргументом командной строки. Например `./auth-proxy -config=./configs/my_app.env`.

Переменные окружения OS перекрывают параметры определенные в конфигурационном файле.
Изменить значение переменной можно в docker-compose.yml

    environment:
      - cookie_name=auth-proxy-dev

## Конфигурационные файлы Postgres

Параметры приложения по умолчанию находятся в в файле  `configs/db.env`.
Файл с параметрами можно указать аргументом командной строки. Например `./auth-proxy -pgconfig=./configs/my_db.env`.

Переменные окружения OS перекрывают параметры определенные в конфигурационном файле.
Изменить значение переменной можно в docker-compose.yml

    environment:
     	- PG_HOST=localhost

## Регистрация нового пользователя

При регистрации нового пользователя (graphql.create_user()) пользователь
автоматически добавляется в группу "users", которая создается автоматически в случае
ее отсутствия.


## Конечные точки

__Приложение в Интернет__

https://auth-proxy.rg.ru/

__Тестовое GUI приложение__ - админка для контроля пользователей, приложений и ролей.

https://auth-proxy.rg.ru/admin/?url=https://auth-proxy.rg.ru


## Конечные точки

__Схема GraphQL__

POST https://auth-proxy.rg.ru/schema

POST https://auth-proxy.rg.ru/graphql


__Captcha__

GET https://auth-proxy.rg.ru/captcha

__Публичный ключ__

GET https://auth-proxy.rg.ru/publickey

__Статистика__

GET https://auth-proxy.rg.ru/stat

__Метрики для Prometheus__

GET https://auth-proxy.rg.ru/metrics

__Список публичных приложений__

GET https://auth-proxy.rg.ru/publicapps

__Тест логирования__

GET https://auth-proxy.rg.ru/logmessage/some_message

Вместо some_message можно вставить собственное сообщение. Сделано для проверки записи логов в Эластик.


__Провайдеры Oauth2__

GET https://auth-proxy.rg.ru/oauthproviders




----------------------------------------------------------


Информация для разработчиков
============================


Запуск приложения
-----------------
    sh/up.sh

Проверка
---------------------

- посылки почты http://localhost:8025/
- Admin Url
    http://localhost:4400/admin/?url=http://localhost:4400




доступ к базе данных pgcli


    pgcli -u pgadmin -d rgru -h 0.0.0.0 -p 5432

В файле ~/.config/pgcli/config определить dsn

    [alias_dsn]
    auth = postgresql://pgadmin:159753@localhost:5432/rgru

Запуск командой

    pgcli -D auth

Пример подключения приложения к auth-proxy
---------------------

 *   [Пример подключения приложения к auth-proxy](./README.newapps.md)



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


Функциональные тесты (End to End) проводятся с помощью  
<https://auth-proxy.rg.ru/admin/?url=https://auth-proxy.rg.ru>


Бенчмарки соединений с БД с пулом и без
--------------------------------------

    go test -run=Bench -benchmem -benchtime=1s -bench=. ./pkg/db

результаты

    Benchmark_local_DB-4           500     3210246 ns/op    18911 B/op     303 allocs/op
    Benchmark_local_DB_pool-4    10000      221768 ns/op     1181 B/op      28 allocs/op
    Benchmark_remote_DB-4          200     7894973 ns/op    22676 B/op     312 allocs/op
    Benchmark_remote_DB_pool-4   10000      210710 ns/op     1181 B/op      28 allocs/op

оптимизация `getKeysAndValues()`

    Benchmark_getKeysAndValues-4              500000              2947 ns/op             848 B/op         28 allocs/op
    Benchmark_getKeysAndValues1-4            1000000               853 ns/op             400 B/op         11 allocs/op



Сборка для фронтэнд разработчиков
----------------------------------------------

    sh/build-frontend-container.sh

## Браузер Chrome и localhost

Если фронтэнд разработчики используют Chrome
им возможно придется установить флаг браузера,
чтобы он допускал самоподписанные сертификаты на localhost.


    chrome://flags/#allow-insecure-localhost



-------------------------------------------------------

Другие команды
--------------------


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

## Деплой на dev

   sh/deploy-dev.sh
    

## Деплой

Перед запуском приложения на сервере должна существовать
сеть `auth_proxy_network`. Сеть порождается командой:

    sudo docker network create auth_proxy_network

Запуск деплоя на продакшн сервере works

    git push origin master:production

или

    sh/deploy.sh



## Контроль деплоя на works

    ssh -i ~/.ssh/deploy_gitupdater_works_open_ssh gitupdater@dockerweb.rgwork.ru


    sudo docker network ls
    sudo docker exec -it auth-proxy-prod bash
    sudo docker logs -f auth-proxy-prod

## Деплой на разные сервера
При деплое на разные сервера может понадобиться использовать разные конфигурационные файлы, для этого в папке configs создается директория (например: gl) с настроечными файлами, характерными для конечного приложения.
Кроме того необходимо создать отдельный build.sh файл (например build-gl.sh) для копирования настроечных файлов на сервер в gitlab-ci.
В самом gitlab деплой разруливается с помощью веток, и шаблонных файлов из директории /gitlab-ci/*.yaml
Ветки:
- production, staging - используются для деплоя на РГ
- production-gl, staging-gl - для деплоя на Год Литературы



## Деплой в ветку staging-gl

    git push origin master:staging-gl

1. Контроль деплоя https://git.rgwork.ru/masterback/auth-proxy/pipelines
2. Можно зайти на dockerweb2 и там проверить

Контроль программы из GUI

https://gl-auth-staging.rg.ru/admin/?url=https://gl-auth-staging.rg.ru#apps

