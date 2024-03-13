# simple_crud_api_go

This is a simple CRUD API written in Go using the [gorilla/mux](https://github.com/gorilla/mux) package

It is based on [this](https://www.youtube.com/watch?v=jFfo23yIWac&t=1234s) tutorial

| Endpoint     | Method |
| ------------ | ------ |
| /movies      | GET    |
| /movies/{id} | GET    |
| /movies      | POST   |
| /movies/{id} | PATCH  |
| /movies/{id} | DELETE |

Here is the JSON representation of a "Movie" object:

```json
"Movie": {
    "id": {
        "description": "Unique app ID",
        "type": "integer"
    },
    "isbn": {
        "description": "Unique movie identifier",
        "type": "string"
    },
    "title": {
        "type": "string"
    },
    "director": {
        "first_name": {
            "type": "string"
        },
        "last_name": {
            "type": "string"
        }
    }
}
```

## Running the program

```console
go run main.go
```
