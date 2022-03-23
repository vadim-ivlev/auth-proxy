## TODO:

- Заменить @ на %40 в урле ссылки в письмах

QUERY

```graphql
query 
list_user_by_usernames(
  $usernames:[String]!,
)
{
  list_user_by_usernames(
  	usernames: $usernames
  ) 
  {
    description
    disabled
    email
    fullname
    id
    username
  }
}
```

VARIABLES

```json
{
  "usernames":["admin", "vadim"]
}
```






https://auth-admin.vercel.app/set-password.html#username=ivlev%40rg.ru&hash=a689d03b-a4f5-4ef9-9b2f-e6138da637ef&authurl=https://gl-auth-staging.rg.ru



https://auth-admin.vercel.app/set-authenticator.html#username=ivlev%40rg.ru&hash=c2cb9146-ec3a-434b-bc9a-35214c22013a&authurl=https://gl-auth-staging.rg.ru


```sql
create view uga as
  select 
  gu.user_email as email, 
  ga.app_appname as appname, 
  ga.rolename as rolename  
  from group_user_role_extended as gu 
  join group_app_role_extended as ga on ga.group_id = gu.group_id
 ;
```