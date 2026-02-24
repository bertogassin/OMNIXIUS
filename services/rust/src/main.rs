// OMNIXIUS — Rust service: video, search, heavy compute
use axum::{routing::get, Json, Router};
use serde::Serialize;
use std::net::SocketAddr;
use tower_http::cors::CorsLayer;

#[derive(Serialize)]
struct Health {
    service: &'static str,
    status: &'static str,
}

async fn health() -> Json<Health> {
    Json(Health {
        service: "omnixius-rust",
        status: "ok",
    })
}

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/health", get(health))
        .layer(CorsLayer::permissive());

    let addr = SocketAddr::from(([0, 0, 0, 0], 8081));
    println!("OMNIXIUS Rust service listening on {}", addr);
    axum::serve(tokio::net::TcpListener::bind(addr).await.unwrap(), app)
        .await
        .unwrap();
}
