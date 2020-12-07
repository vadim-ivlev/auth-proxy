## TODO:


- Простой механизм. Считывать конфигурационные файлы пока не произойдет успешное соединение с базой данных.
  Это нужно для того, чтобы Лебедев мог изменять конфигурационные файлы после старта программы.
- Сложный механизм. Следить за изменением конфигурационных файлов. Если они изменились Переодсоединяться к базе данных.
- Исправить баг с передачей списков ролей конечному приложению.


Add PUT, DELETE verbs
add auth_proxy_network



- [GIN] 2019/11/2   8 - 18:27:49 | 404 |      14.536µs |             ::1 | GET      /bundles/jurist/fonts/rg/noto.woff.min.css




- Oauth2 (finish logout for yandex, vk, mail.ru)




## development
- get rid of gin using <https://github.com/abbot/go-http-auth>
- get rid of sessions
    https://github.com/Depado/gin-auth-example/blob/master/main.go


## optimize
- (postponed) change search Like statement from logical OR to concatenation of fields ||
- add DBStats function to GraphQL


