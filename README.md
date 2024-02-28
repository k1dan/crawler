# crawler
Repository contains crawler for ebay.com. To run application define envs in .env file and run. </br>
``CONDITIONS_TO_PARSE ``env specify conditions of products that should be parsed. </br>
Several conditions can be specified separated by comma (example: "new, used"). </br>
If empty string then no conditions are applied, by default env is empty string
``
    go run cmd/app/main.go
``

## envs:
```
    APP_NAME=
    LOG_LEVEL=
    FILE_SAVE_PATH=
    PARSER_WORKER_AMOUNT=
    FILE_SAVER_WORKER_AMOUNT=
    CONDITIONS_TO_PARSE=
    PARSING_START_URL=
```

## Future improvement steps
- Refactor workers
- Make parsing constructor
- Add Tests
- Add Containerization 
