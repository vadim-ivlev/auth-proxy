auth-proxy
---------

https://auth-proxy.rg.ru

### auth-proxy authorization service

Права:
API Яндекс.Паспорта
Доступ к логину, имени и фамилии, полу

Доступ к адресу электронной почты

Доступ к дате рождения

Доступ к портрету пользователя

-----


ID: 39cd5970ac5148ca901f7298d3f73246

Пароль: d0061eaba5f244ecb1fd10e9069a418f

Callback URL: http://localhost:4400/oauth2callback



```
https://oauth.yandex.ru/authorize?
   response_type=code
 & client_id=<идентификатор приложения>
[& device_id=<идентификатор устройства>]
[& device_name=<имя устройства>]
[& redirect_uri=<адрес перенаправления>]
[& login_hint=<имя пользователя или электронный адрес>]
[& scope=<запрашиваемые необходимые права>]
[& optional_scope=<запрашиваемые опциональные права>]
[& force_confirm=yes]
[& state=<произвольная строка>]

```

Example: 

https://oauth.yandex.ru/authorize?response_type=code&client_id=39cd5970ac5148ca901f7298d3f73246&redirect_uri=http://localhost:4400/oauth2callback&state123
