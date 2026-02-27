// OMNIXIUS — Rust service (stack 1): video, search, heavy compute
use axum::{extract::Json, routing::get, routing::post, Router};
use serde::{Deserialize, Serialize};
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

#[derive(Deserialize)]
struct RankRequest {
    #[serde(default)]
    items: Vec<serde_json::Value>,
    #[serde(default)]
    sort: String,
}

#[derive(Serialize)]
struct RankResponse {
    items: Vec<serde_json::Value>,
}

async fn rank(Json(req): Json<RankRequest>) -> Json<RankResponse> {
    let mut items = req.items;
    if req.sort == "rating" {
        items.sort_by(|a, b| {
            let ra = a.get("rating_avg").and_then(|v| v.as_f64()).unwrap_or(0.0);
            let rb = b.get("rating_avg").and_then(|v| v.as_f64()).unwrap_or(0.0);
            ra.partial_cmp(&rb).unwrap_or(std::cmp::Ordering::Equal).reverse()
        });
    } else if req.sort == "distance" {
        items.sort_by(|a, b| {
            let da = a.get("distance_km").and_then(|v| v.as_f64()).unwrap_or(f64::MAX);
            let db = b.get("distance_km").and_then(|v| v.as_f64()).unwrap_or(f64::MAX);
            da.partial_cmp(&db).unwrap_or(std::cmp::Ordering::Equal)
        });
    }
    Json(RankResponse { items })
}

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/health", get(health))
        .route("/rank", post(rank))
        .layer(CorsLayer::permissive());

    let addr = SocketAddr::from(([0, 0, 0, 0], 8081));
    println!("OMNIXIUS Rust service (stack 1) listening on {}", addr);
    axum::serve(tokio::net::TcpListener::bind(addr).await.unwrap(), app)
        .await
        .unwrap();
}
