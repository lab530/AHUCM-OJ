use std::ops::Deref;

use actix_web::{App, HttpServer};
use api::{
    internal::{ping_get, ping_post, reload_config},
    judge::submit,
};
use constants::SERVER_ADDR;

mod api;
mod constants;
mod model;
mod service;
mod util;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init();
    let server_addr = SERVER_ADDR.deref();
    log::debug!("will bind on {server_addr}");
    HttpServer::new(|| {
        App::new()
            .service(ping_get)
            .service(ping_post)
            .service(submit)
            .service(reload_config)
    })
    .bind(server_addr)?
    .run()
    .await
}
