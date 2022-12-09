## TODO:


 .dockerignore                      |   3 +-
 .helm/templates/90-ingress.yaml    |  17 +-------
 .helm/values.yaml                  |  12 +++---
 README.md                          | 145 +++++++++++++++++++++++++++++++++---------------------------------
 admin/mustache/app.html            |  90 +++++++++++++++++++++++------------------
 docker-compose-dev.yml             |  12 +++---
 docker-compose-front.yml           |  26 ++++++------
 migrations/01_create_tables.up.sql |  78 +++++++++++++++---------------------
 migrations/04_add_data.up.sql      | 153 ++++++++++++++++++++++++++++++++----------------------------------



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


# Working with git via RG

1. Pull the latest changes to rg

```
ssh rg 'cd auth-proxy && git pull origin master'

```


2. 
On the local computer copy from rg

```
rsync -avh --delete rg:/home/ivlev/auth-proxy ~/projects/
```

3. Do the work and commit the changes.

Copy to rg
```
rsync -avh --delete  ~/projects/auth-proxy rg:/home/ivlev/
```

