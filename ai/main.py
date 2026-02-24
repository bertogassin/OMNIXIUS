"""
OMNIXIUS — AI: neural nets, recognition (Python).
Root of our own AI — start here; later models will be strongest here.
"""
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


class ChatOut(BaseModel):
    reply: str
    model: str = "omnixius-ai-v0"


@app.get("/health")
def health():
    return {"service": "omnixius-ai", "status": "ok"}


@app.get("/")
def root():
    return {"message": "OMNIXIUS AI — neural nets, recognition. Root of our AI."}


@app.post("/chat", response_model=ChatOut)
def chat(body: ChatIn):
    """Root of OMNIXIUS AI. Placeholder: echo + intent. Replace with real model later."""
    msg = (body.message or "").strip()
    if not msg:
        return ChatOut(reply="Send a message to start. This is the root of OMNIXIUS AI — our models will run here.")
    # Placeholder reply; later: inference from our own model
    reply = f"You said: {msg[:200]}. This is OMNIXIUS AI root — full model coming next."
    return ChatOut(reply=reply)
