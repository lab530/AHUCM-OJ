use actix_web::{get, post, HttpResponse, Responder};

use crate::util::config::GLOB_CONFIG;

#[get("/api/v1/internal/ping")]
pub async fn ping_get() -> impl Responder {
    HttpResponse::Ok().body("pong")
}

#[post("/api/v1/internal/ping")]
pub async fn ping_post() -> impl Responder {
    HttpResponse::Ok().body("pong")
}

#[post("/api/v1/internal/config/reload")]
pub async fn reload_config() -> impl Responder {
    if let Some(e) = GLOB_CONFIG.lock().unwrap().reload() {
        HttpResponse::InternalServerError().body(e.clone())
    } else {
        HttpResponse::Ok().finish()
    }
}
