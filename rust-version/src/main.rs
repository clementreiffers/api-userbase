use actix_web::web::Data;
use actix_web::{get, post, web, App, HttpResponse, HttpServer};
use futures::stream::TryStreamExt;
use mongodb::bson::doc;
use mongodb::results::InsertOneResult;
use mongodb::Database;
use mongodb::{options::ClientOptions, Client};
use serde::{Deserialize, Serialize};
use serde_json::json;

#[derive(Debug, Serialize, Deserialize)]
struct User {
    name: String,
}

#[get("/ping")]
async fn greet() -> HttpResponse {
    HttpResponse::Ok().body(json!({"message":"pong"}).to_string())
}

async fn find_user(database: Data<mongodb::Database>, user: User) -> mongodb::Cursor<User> {
    let query = doc! {"name":user.name};

    database
        .collection("users")
        .find(query, None)
        .await
        .ok()
        .expect("unable to find useri named {user.name}")
}

#[get("/get-user")]
async fn get_user(payload: web::Json<User>, database: Data<mongodb::Database>) -> HttpResponse {
    let cursor: Vec<User> = find_user(database, payload.into_inner())
        .await
        .try_collect()
        .await
        .unwrap();
    let data = serde_json::to_string(&cursor).expect("failed to serialize data");
    println!("{:?}", data);
    HttpResponse::Ok().body(data)
}

async fn insert_user(database: &Data<Database>, user: User) -> InsertOneResult {
    database
        .collection::<User>("users")
        .insert_one(user, None)
        .await
        .expect("unable to insert a user")
}

#[post("/add-user")]
async fn add_user(payload: web::Json<User>, database: Data<mongodb::Database>) -> HttpResponse {
    let insert_user_result = insert_user(&database, payload.into_inner()).await;
    let inserted_id = format!("{:?}", insert_user_result.inserted_id);

    if inserted_id != "" {
        HttpResponse::Ok().body(json!({"message":"sucessfully uploaded!"}).to_string())
    } else {
        HttpResponse::InternalServerError()
            .body(json!({"message":"cannot upload this user"}).to_string())
    }
}

async fn connect_mongo_database(uri: &str) -> Database {
    let client_options = ClientOptions::parse(uri)
        .await
        .expect("unable to connect to database uri");
    let client = Client::with_options(client_options).expect("unable to create a mongo client");
    let database = client.database("userBase");
    database
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    std::env::set_var("RUST_LOG", "debug");
    env_logger::init();
    let uri = "mongodb://localhost:27017/userBase";

    let database = connect_mongo_database(uri).await;
    HttpServer::new(move || {
        App::new()
            .app_data(Data::new(database.clone()))
            .service(greet)
            .service(get_user)
            .service(add_user)
    })
    .bind(("127.0.0.1", 8080))?
    .run()
    .await
}
