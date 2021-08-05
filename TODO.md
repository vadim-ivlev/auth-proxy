## TODO:

- When proxying to url with path http://abc.com/a/b/c. Remove prefix?
- Do routing by hands `apps.Any("/:appname/*proxypath", Proxy) -> apps.Any("*proxypath", Proxy)`


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