# langchain-go-example
## How to run
1. set environment variables
```bash
export OPENAI_API_KEY=YOUR_OPENAI_API_KEY
````

2. Install docker and docker-compose
```bash
docker-compose up
```

3. Migrate VectorStores
```
go run script/weaviate/main.go
```

4. Insert FAQ Data
```
go run script/insert_docs_csv/main.go
go run script/insert_docs_html/main.go
```

5. Run the main program
```
go run main.go
```