# Token Store API

Simple api for bulk storing tokens and retrieving (and deleting) single tokens.
Uses sqlite as the default db.

## API
`POST /tokens`

Store tokens in the store to be fetched later

Request Payload:
```
{"tokens":["a","b","c"]}
```


`GET /tokens`

Fetch a single token from the database and remove it.

Response Body:
```
{"token": "a"}
```