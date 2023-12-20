# rust-version

## build it 

to build it in release mode, run `cargo build --release`

### run it 

you first need to setup a `mongodb` uri, so launch it locally running `docker run -p 27017:27017 -d mongo` or setup it on the [mongodb website](www.mongodb.com)

to run it without building it in release mode, launch `cargo run -- --mongo-uri mongodb://localhost:27017/userBase`

to run it while the build is finised, run `./target/release/api-userbase --mongo-uri mongodb://localhost:27017/userBase` 
