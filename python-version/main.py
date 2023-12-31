from fastapi import FastAPI
import pymongo
from pydantic import BaseModel
import argparse

parser = argparse.ArgumentParser()
parser.add_argument("-h", "--hostname", help="mongodb host")

args = parser.parse_args()

mongo_uri = args.hostname
database_name = "userBase"
collection_name = "users"

collection = pymongo.MongoClient(mongo_uri)[database_name][collection_name]
app = FastAPI()


class User(BaseModel):
    name: str


@app.get("/get-user")
async def get_user(user: User):
    return str(list(map(lambda document: document, collection.find(user.model_dump()))))


@app.post("/add-user")
async def add_user(user: User):
    insert_result = collection.insert_one(user.model_dump())
    return {
        "message": "successfully uploaded"
        if insert_result.inserted_id
        else "internal error"
    }


@app.get("/ping")
async def pong():
    return {"message": "pong"}
