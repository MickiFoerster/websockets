use std::{net::SocketAddr, path::PathBuf};

use axum::extract::connect_info::ConnectInfo;
use axum::{
    extract::ws::{Message, WebSocket, WebSocketUpgrade},
    response::IntoResponse,
    routing::get,
    Router,
};
use axum_extra::TypedHeader;
use futures::{sink::SinkExt, stream::StreamExt};
use serde_json::json;
use tower_http::{
    services::ServeDir,
    trace::{DefaultMakeSpan, TraceLayer},
};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

#[tokio::main]
async fn main() {
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env().unwrap_or_else(|_| {
                "form_validation_websockets=debug,tower_http=info,tower=info".into()
            }),
        )
        .with(tracing_subscriber::fmt::layer())
        .init();

    let assets_dir = PathBuf::from(".").join("assets");

    let app = Router::new()
        .fallback_service(ServeDir::new(assets_dir).append_index_html_on_directories(true))
        .route("/ws", get(ws_handler))
        .layer(
            TraceLayer::new_for_http()
                .make_span_with(DefaultMakeSpan::default().include_headers(true)),
        );

    let listener = tokio::net::TcpListener::bind("127.0.0.1:3000")
        .await
        .unwrap();
    tracing::info!("listening on {}", listener.local_addr().unwrap());
    axum::serve(
        listener,
        app.into_make_service_with_connect_info::<SocketAddr>(),
    )
    .await
    .unwrap();
}

async fn ws_handler(
    ws: WebSocketUpgrade,
    user_agent: Option<TypedHeader<headers::UserAgent>>,
    ConnectInfo(addr): ConnectInfo<SocketAddr>,
) -> impl IntoResponse {
    let user_agent = if let Some(TypedHeader(user_agent)) = user_agent {
        user_agent.to_string()
    } else {
        String::from("Unknown browser")
    };
    tracing::info!("`{user_agent}` at {addr} connected.");

    ws.on_upgrade(move |socket| handle_socket(socket, addr))
}

async fn handle_socket(socket: WebSocket, who: SocketAddr) {
    let (tx, mut rx) = tokio::sync::mpsc::channel(10);

    let (mut sender, mut receiver) = socket.split();

    // task for sending messages to the client
    let mut send_task = tokio::spawn(async move {
        tracing::info!("sender waits for message");
        while let Some(msg) = rx.recv().await {
            tracing::info!("sender received {msg:?} via internal channel");
            if sender.send(Message::Text(msg)).await.is_err() {
                return;
            }
        }
        tracing::info!("sender finished");
    });

    // receiver task
    let mut recv_task = tokio::spawn(async move {
        while let Some(Ok(msg)) = receiver.next().await {
            let msg = match msg {
                Message::Text(t) => t,
                Message::Binary(_) => todo!(),
                Message::Ping(_) => todo!(),
                Message::Pong(_) => todo!(),
                Message::Close(_) => todo!(),
            };

            tracing::info!("received: {:#?}", msg);

            match tx
                .send(
                    json!(
                    {
                        "success": "true",
                        "message": format!("{}", msg)
                    })
                    .to_string(),
                )
                .await
            {
                Ok(v) => tracing::debug!("Ok: {v:?}"),
                Err(v) => tracing::debug!("Err: {v:?}"),
            }
            tracing::info!("value send from receiver to sender");
        }
    });

    // If any one of the tasks exit, abort the other.
    tokio::select! {
        rv_a = (&mut send_task) => {
            match rv_a {
                Ok(()) => println!("{who} connection closed"),
                Err(e) => println!("Error sending messages {e:?}")
            }
            recv_task.abort();
        },
        rv_b = (&mut recv_task) => {
            match rv_b {
                Ok(()) => println!("connection receiving part closed"),
                Err(e) => println!("Error receiving messages {e:?}")
            }
            send_task.abort();
        }
    }

    // returning from the handler closes the websocket connection
    println!("Websocket context {who} destroyed");
}