# crawler
Repository contains crawler for ebay.com. To run application define envs in .env file and run 
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
