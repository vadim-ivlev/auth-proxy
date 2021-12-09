## TODO:

this is a classic branch

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