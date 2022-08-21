### SpaceTrouble

To start:
   ```
   docker-compose up
   ```
To shutdown: 
```
   docker-compose down
```

app will be available on http://localhost:8000

---------------------------------------------------------

### Available endpoints:

#### List of destinations:
```curl
curl --request GET 'http://127.0.0.1:8000/api/v1/destinations'
```

response: 
```json
[
    {
        "id": "1",
        "name": "Mars"
    },
    {
        "id": "2",
        "name": "IO"
    }
]
```

#### Order 

```curl 
curl --request POST 'http://127.0.0.1:8000/api/v1/orders' \
--header 'Content-Type: application/json' \
--data-raw '{
    "first_name": "Vasyl",
    "last_name": "Osypchuk",
    "gender": "male",
    "birthday_year": 2000,
    "birthday_month": 3,
    "birthday_day": 1,
    "launchpad_id": "5e9e4501f509094ba4566f84",
    "destination_id": "1",
    "launch_date": "2022-09-04T00:00:00-07:00"
}'
```

Success response:
```json
{
    "id": "e531b91b-46b6-44d0-937c-226c7cb51bb8"
}
```

Possible error codes:<br>
   <strong>400</strong> - invalid data  (like missing fields, launch date in the past, launchpad or destination is not exists)
   <strong>406</strong> - launchpad or busy or has another destination for provided launch date

#### List of orders

```curl
curl --request GET 'http://127.0.0.1:8000/api/v1/orders?limit=10&offset=0'
```

response:
```json
{
    "docs": [
        {
            "id": "e531b91b-46b6-44d0-937c-226c7cb51bb8",
            "first_name": "Vasyl",
            "last_name": "Osypchuk",
            "gender": "male",
            "birthday_year": 2000,
            "birthday_month": 3,
            "birthday_day": 1,
            "launchpad_id": "5e9e4501f509094ba4566f84",
            "destination_id": "1",
            "launch_date": "2022-09-04T07:00:00Z",
            "created_at": "2022-08-21T12:57:26.950964Z"
        }
    ],
    "limit": 10,
    "offset": 0
}
```

#### Get order by id

```curl
curl --request GET 'http://127.0.0.1:8000/api/v1/orders/{id}'
```
response:
```json
{
    "id": "e531b91b-46b6-44d0-937c-226c7cb51bb8",
    "first_name": "Vasyl",
    "last_name": "Osypchuk",
    "gender": "male",
    "birthday_year": 2000,
    "birthday_month": 3,
    "birthday_day": 1,
    "launchpad_id": "5e9e4501f509094ba4566f84",
    "destination_id": "1",
    "launch_date": "2022-09-04T07:00:00Z",
    "created_at": "2022-08-21T12:57:26.950964Z"
}
```

#### Delete order
```curl
curl --request DELETE 'http://127.0.0.1:8000/api/v1/orders/{id}'
```

returns 204 without content


