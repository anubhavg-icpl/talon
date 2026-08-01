"""OpenAI-compatible SLM server for Talon.

Endpoints:
  GET  /health              — liveness + engine info
  GET  /v1/models           — model list
  POST /v1/chat/completions — chat (stream=true → SSE, OpenAI chunk format)

Go talks to this exactly like any OpenAI-compatible backend (LLM_PROVIDER=onnx).
"""

from __future__ import annotations

import json
import logging
import os
import time
import uuid
from contextlib import asynccontextmanager
from typing import Any, AsyncGenerator, Optional

from fastapi import FastAPI, HTTPException, Request
from fastapi.responses import JSONResponse, StreamingResponse
from pydantic import BaseModel, Field

from .engine import (
    ChatMessage,
    GenerateRequest,
    engine_ready,
    get_engine,
    list_models,
)

logging.basicConfig(
    level=os.getenv("LOG_LEVEL", "INFO"),
    format="%(asctime)s %(levelname)s %(name)s: %(message)s",
)
logger = logging.getLogger("onnx-slm")

PORT = int(os.getenv("SLM_PORT", "8090"))
# Eager load on startup so first request is not cold (unless SLM_LAZY=1).
LAZY = os.getenv("SLM_LAZY", "0") == "1"


@asynccontextmanager
async def lifespan(_app: FastAPI):
    if not LAZY:
        try:
            eng = get_engine()
            logger.info("engine ready: %s", eng.info())
        except Exception:
            logger.exception("engine preload failed — will retry on first request")
    yield


app = FastAPI(
    title="Talon SLM Runtime",
    description="Local SmolLM / ONNX Runtime OpenAI-compatible inference",
    version="1.0.0",
    lifespan=lifespan,
)


class ChatMessageIn(BaseModel):
    role: str
    content: str | list[Any] = ""


class ChatCompletionsRequest(BaseModel):
    model: str = "smollm"
    messages: list[ChatMessageIn]
    max_tokens: Optional[int] = Field(default=None, alias="max_tokens")
    max_completion_tokens: Optional[int] = None
    temperature: Optional[float] = 0.3
    stream: bool = False
    # Accepted and ignored (OpenAI client may send them).
    tools: Optional[list[Any]] = None
    tool_choice: Optional[Any] = None
    top_p: Optional[float] = None
    stop: Optional[Any] = None
    n: Optional[int] = 1
    user: Optional[str] = None

    model_config = {"populate_by_name": True, "extra": "ignore"}


def _content_text(content: str | list[Any]) -> str:
    if isinstance(content, str):
        return content
    # Multimodal-style content parts: take text pieces only.
    parts: list[str] = []
    for p in content:
        if isinstance(p, dict) and p.get("type") == "text":
            parts.append(str(p.get("text", "")))
        elif isinstance(p, str):
            parts.append(p)
    return "".join(parts)


def _to_messages(msgs: list[ChatMessageIn]) -> list[ChatMessage]:
    out: list[ChatMessage] = []
    for m in msgs:
        role = m.role if m.role in ("system", "user", "assistant", "tool") else "user"
        out.append(ChatMessage(role=role, content=_content_text(m.content)))
    return out


def _completion_id() -> str:
    return f"chatcmpl-slm-{uuid.uuid4().hex[:12]}"


@app.get("/health")
def health():
    try:
        eng = get_engine()
        info = eng.info()
        return {
            "status": "ok",
            "backend": info.backend,
            "model_id": info.model_id,
            "device": info.device,
            "ready": info.ready,
            "load_ms": info.load_ms,
            "extra": info.extra,
        }
    except Exception as e:
        return JSONResponse(
            status_code=503,
            content={"status": "error", "ready": False, "detail": str(e)},
        )


@app.get("/v1/models")
def models():
    return {"object": "list", "data": list_models()}


@app.get("/models")
def models_alias():
    return models()


def _sse_chunk(cid: str, model: str, delta: dict, finish: Optional[str] = None) -> str:
    body = {
        "id": cid,
        "object": "chat.completion.chunk",
        "created": int(time.time()),
        "model": model,
        "choices": [
            {
                "index": 0,
                "delta": delta,
                "finish_reason": finish,
            }
        ],
    }
    return f"data: {json.dumps(body, ensure_ascii=False)}\n\n"


def _stream_response(req: ChatCompletionsRequest) -> StreamingResponse:
    eng = get_engine()
    cid = _completion_id()
    model = req.model or eng.info().model_id
    max_tokens = req.max_tokens or req.max_completion_tokens or 256
    gen_req = GenerateRequest(
        messages=_to_messages(req.messages),
        max_tokens=max_tokens,
        temperature=req.temperature if req.temperature is not None else 0.3,
        model=model,
        stream=True,
    )

    def event_gen():
        # Role header (OpenAI clients expect this first).
        yield _sse_chunk(cid, model, {"role": "assistant", "content": ""})
        t0 = time.perf_counter()
        n = 0
        try:
            for piece in eng.stream(gen_req):
                n += 1
                yield _sse_chunk(cid, model, {"content": piece})
        except Exception as e:
            logger.exception("stream failed")
            err = {"error": {"message": str(e), "type": "server_error"}}
            yield f"data: {json.dumps(err)}\n\n"
            yield "data: [DONE]\n\n"
            return
        yield _sse_chunk(cid, model, {}, finish="stop")
        yield "data: [DONE]\n\n"
        ms = int((time.perf_counter() - t0) * 1000)
        logger.info("streamed %d pieces in %dms model=%s", n, ms, model)

    return StreamingResponse(
        event_gen(),
        media_type="text/event-stream",
        headers={
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
            "X-Accel-Buffering": "no",
        },
    )


@app.post("/v1/chat/completions")
def chat_completions(req: ChatCompletionsRequest):
    if not req.messages:
        raise HTTPException(status_code=400, detail="messages required")
    # Tool-calling is not supported on SLMs here; refuse clearly so Go can
    # fall back to a tool-capable provider for agent turns.
    if req.tools:
        raise HTTPException(
            status_code=400,
            detail="tools/function-calling not supported on the SLM runtime; "
            "use LLM_PROVIDER=openai|bedrock|ollama for agent tool loops, "
            "or omit tools for plain chat streaming",
        )
    try:
        if req.stream:
            return _stream_response(req)
        eng = get_engine()
        max_tokens = req.max_tokens or req.max_completion_tokens or 256
        text = eng.generate(
            GenerateRequest(
                messages=_to_messages(req.messages),
                max_tokens=max_tokens,
                temperature=req.temperature if req.temperature is not None else 0.3,
                model=req.model,
            )
        )
        cid = _completion_id()
        model = req.model or eng.info().model_id
        return {
            "id": cid,
            "object": "chat.completion",
            "created": int(time.time()),
            "model": model,
            "choices": [
                {
                    "index": 0,
                    "message": {"role": "assistant", "content": text},
                    "finish_reason": "stop",
                }
            ],
            "usage": {
                "prompt_tokens": 0,
                "completion_tokens": 0,
                "total_tokens": 0,
            },
        }
    except HTTPException:
        raise
    except Exception as e:
        logger.exception("chat failed")
        raise HTTPException(status_code=500, detail=str(e)) from e


# Alias without /v1 prefix (some clients omit it).
@app.post("/chat/completions")
def chat_completions_alias(req: ChatCompletionsRequest):
    return chat_completions(req)


@app.get("/")
def root():
    return {
        "service": "talon-onnx-slm",
        "docs": "/docs",
        "health": "/health",
        "openai": "/v1/chat/completions",
        "ready": engine_ready(),
    }


def main() -> None:
    import uvicorn

    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=PORT,
        log_level=os.getenv("LOG_LEVEL", "info").lower(),
        # Workers=1: model lives in one process (multi-worker would multiply RAM).
        workers=1,
        timeout_keep_alive=75,
    )


if __name__ == "__main__":
    main()
