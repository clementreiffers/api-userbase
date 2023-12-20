from fastapi import FastAPI
import pymongo

app = FastAPI()

mongo_uri = "mongodb://localhost:27017/userBase"

client = pymongo.MongoClient(mongo_uri)
database_name = "userBase"
collection_name = "users"
collection = client[database_name][collection_name]


@app.get("/get-user")
async def get_user(json_body: dict):
    return str(list(map(lambda document: document, collection.find(json_body))))


@app.post("/add-user")
async def add_user(json_body: dict):
    insert_result = collection.insert_one(json_body)

    return {
        "message": "successfully uploaded"
        if insert_result.inserted_id
        else "internal error"
    }


@app.get("/ping")
async def pong():
    return {"message": "pong"}
