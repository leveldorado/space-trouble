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

### Order 

