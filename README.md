# llm-qa-go-example
## How to run
1. Install docker and docker-compose
```bash
docker-compose up
```

2. Migrate VectorStores
```
go run script/weaviate/main.go
```

3. Insert FAQ Data
```
go run script/insert_docs/main.go
```

4. Run the main program
```
go run main.go
```