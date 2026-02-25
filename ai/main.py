"""
OMNIXIUS — AI: neural nets, recognition (Python).
Root of our own AI — start here; later models will be strongest here.
A1: my orders. A2: create order, conversations summary. See ai/README.md for intent catalog.
"""
import re
from typing import Optional, Tuple

import httpx
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel

app = FastAPI(title="OMNIXIUS AI")
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)


class ChatIn(BaseModel):
    message: str
    module: Optional[str] = None
    api_token: Optional[str] = None
    api_url: Optional[str] = None


class ChatOut(BaseModel):
    reply: str
    model: str = "omnixius-ai-v0"


def _normalize(s: str) -> str:
    return (s or "").lower().strip()


def _is_my_orders_intent(msg: str) -> bool:
    """Detect if user wants to see their orders (EN/RU/FR)."""
    n = _normalize(msg)
    if not n or len(n) < 2:
        return False
    # English
    if re.search(r"\b(my\s+)?orders?\b", n) or "show orders" in n or "list orders" in n:
        return True
    # Russian
    if re.search(r"заказ|мои заказы|покажи заказы|список заказов", n):
        return True
    # French
    if re.search(r"mes\s+commandes|commandes?|liste des commandes", n):
        return True
    return False


def _fetch_my_orders(api_url: str, token: str) -> str:
    """Call backend-go GET /api/orders/my; return formatted text or error."""
    base = (api_url or "").rstrip("/")
    if not base:
        return "Backend URL is not set. Set it in the app (e.g. login page) and try again."
    url = f"{base}/api/orders/my"
    try:
        with httpx.Client(timeout=10.0) as client:
            r = client.get(
                url,
                headers={"Authorization": f"Bearer {token}", "Content-Type": "application/json"},
            )
        if r.status_code == 401:
            return "You are not signed in or the session expired. Please sign in again in the app."
        if r.status_code != 200:
            return f"Backend returned {r.status_code}. Try again later."
        data = r.json()
        # Backend may return a list (current) or { asBuyer, asSeller } (if extended)
        orders = []
        if isinstance(data, list):
            orders = data[:30]
        else:
            orders = (data.get("asBuyer") or []) + (data.get("asSeller") or [])
            orders = orders[:30]
        lines = []
        for o in orders:
            title = o.get("title") or "—"
            status = o.get("status") or "—"
            price = o.get("price")
            price_s = f" — {price}" if price is not None else ""
            lines.append(f"• {title} — {status}{price_s}")
        if not lines:
            return "You have no orders yet."
        return "\n".join(lines)
    except httpx.RequestError as e:
        return f"Cannot reach the backend. Check that it is running at {base} and try again."
    except Exception as e:
        return f"Something went wrong: {str(e)[:200]}"


def _is_create_order_intent(msg: str) -> Tuple[bool, Optional[int]]:
    """Detect create-order intent and extract product_id. Returns (is_intent, product_id or None)."""
    n = _normalize(msg)
    if not n or len(n) < 3:
        return False, None
    if _is_my_orders_intent(msg):
        return False, None
    # Verb: create order, order product, buy, оформить заказ, заказать, commander
    has_verb = bool(
        re.search(r"(?:create|place|make)\s+(?:an?\s+)?order|order\s+(?:product\s+)?", n)
        or re.search(r"buy\s+(?:product\s+)?", n)
        or re.search(r"оформить\s+заказ|заказать\s+(?:продукт\s+)?", n)
        or re.search(r"commander\s+(?:produit\s+)?", n)
    )
    if not has_verb:
        return False, None
    m = re.search(r"(\d+)", n)
    if m:
        pid = int(m.group(1))
        if 1 <= pid <= 999999:
            return True, pid
    return True, None  # verb present but no id → ask for product id


def _do_create_order(api_url: str, token: str, product_id: int) -> str:
    """POST /api/orders with product_id. Return success or error text."""
    base = (api_url or "").rstrip("/")
    if not base:
        return "Backend URL is not set."
    url = f"{base}/api/orders"
    try:
        with httpx.Client(timeout=10.0) as client:
            r = client.post(
                url,
                headers={"Authorization": f"Bearer {token}", "Content-Type": "application/json"},
                json={"product_id": product_id},
            )
        if r.status_code == 401:
            return "Please sign in again."
        if r.status_code == 400:
            data = r.json() if r.text else {}
            err = data.get("error", r.text or "Bad request")
            return f"Could not create order: {err}"
        if r.status_code == 404:
            return f"Product {product_id} not found."
        if r.status_code != 200:
            return f"Backend returned {r.status_code}."
        data = r.json()
        oid = data.get("id", "—")
        return f"Order created. Order ID: {oid}. You can check status in Orders."
    except httpx.RequestError:
        return "Cannot reach the backend. Check that it is running."
    except Exception as e:
        return f"Error: {str(e)[:150]}"


def _is_conversations_summary_intent(msg: str) -> bool:
    """Detect intent: summary of messages / conversations (EN/RU/FR)."""
    n = _normalize(msg)
    if not n or len(n) < 2:
        return False
    if re.search(r"\b(my\s+)?(messages?|conversations?|chats?|mail|inbox)\b", n):
        return True
    if re.search(r"summary|сводка|кратко|résumé|résume", n) and re.search(r"message|письм|переписк|conversation", n):
        return True
    if re.search(r"письма|переписки|сообщения|диалоги", n):
        return True
    if re.search(r"mes\s+messages?|conversations?", n):
        return True
    return False


def _fetch_conversations_summary(api_url: str, token: str) -> str:
    """GET /api/conversations; return short summary (last_message, other, unread)."""
    base = (api_url or "").rstrip("/")
    if not base:
        return "Backend URL is not set."
    url = f"{base}/api/conversations"
    try:
        with httpx.Client(timeout=10.0) as client:
            r = client.get(
                url,
                headers={"Authorization": f"Bearer {token}", "Content-Type": "application/json"},
            )
        if r.status_code == 401:
            return "Please sign in again."
        if r.status_code != 200:
            return f"Backend returned {r.status_code}."
        data = r.json()
        if not data:
            return "No conversations yet."
        lines = []
        for c in data[:15]:
            other = (c.get("other") or {}).get("name") or (c.get("other") or {}).get("email") or "—"
            last = (c.get("last_message") or "")[:80]
            unread_s = " (unread)" if c.get("unread") else ""
            lines.append(f"• {other}: {last or '—'}{unread_s}")
        return "\n".join(lines) if lines else "No conversations yet."
    except httpx.RequestError:
        return "Cannot reach the backend."
    except Exception as e:
        return f"Error: {str(e)[:150]}"


@app.get("/health")
def health():
    return {"service": "omnixius-ai", "status": "ok"}


@app.get("/")
def root():
    return {"message": "OMNIXIUS AI — neural nets, recognition. Root of our AI."}


@app.post("/chat", response_model=ChatOut)
def chat(body: ChatIn):
    """Root of OMNIXIUS AI. A1/A2: intents → backend-go (orders, create order, conversations)."""
    msg = (body.message or "").strip()
    if not msg:
        return ChatOut(reply="Send a message to start. This is the root of OMNIXIUS AI — our models will run here.")

    token = (body.api_token or "").strip()
    api_url = (body.api_url or "").strip()
    need_auth = "Sign in and set the backend API URL (e.g. on the login page) so I can act on your behalf."

    # A1: my orders
    if _is_my_orders_intent(msg):
        if token and api_url:
            return ChatOut(reply=_fetch_my_orders(api_url, token))
        return ChatOut(reply=need_auth)

    # A2: create order (product_id extracted from message)
    is_create, product_id = _is_create_order_intent(msg)
    if is_create and product_id and token and api_url:
        return ChatOut(reply=_do_create_order(api_url, token, product_id))
    if is_create and (not token or not api_url):
        return ChatOut(reply=need_auth)
    if is_create and not product_id:
        return ChatOut(reply="Specify the product ID, e.g. «create order 5» or «order product 5».")

    # A2: conversations / messages summary
    if _is_conversations_summary_intent(msg) and token and api_url:
        return ChatOut(reply=_fetch_conversations_summary(api_url, token))
    if _is_conversations_summary_intent(msg) and (not token or not api_url):
        return ChatOut(reply=need_auth)

    # Fallback: hint at available actions
    reply = (
        f"You said: {msg[:200]}. "
        "I can: show **my orders**, **create order** (e.g. «order product 5»), or give a **summary of messages**. "
        "Say one of these when signed in."
    )
    return ChatOut(reply=reply)
