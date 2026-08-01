"""Inference backends for local SLMs.

Primary path: HuggingFace transformers (SmolLM2-135M/360M) — low-latency CPU
streaming with chat templates. Optional path: ONNX Runtime when a model.onnx
(or model_q4.onnx) is present under MODEL_DIR.

Tokens are yielded as soon as they decode so the Go API can SSE them to the
dashboard in milliseconds.
"""

from __future__ import annotations

import logging
import os
import threading
import time
from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from pathlib import Path
from typing import Generator, Iterable, Optional

logger = logging.getLogger("onnx-slm.engine")

DEFAULT_HF_ID = os.getenv("SLM_MODEL_ID", "HuggingFaceTB/SmolLM2-360M-Instruct")
MODEL_DIR = Path(os.getenv("MODEL_DIR", "/models/smollm"))
MAX_NEW_TOKENS = int(os.getenv("SLM_MAX_NEW_TOKENS", "256"))
DEFAULT_TEMPERATURE = float(os.getenv("SLM_TEMPERATURE", "0.3"))
# Prefer smaller model on low-RAM hosts when SLM_MODEL_ID is unset and 360M fails.
FALLBACK_HF_ID = os.getenv("SLM_FALLBACK_MODEL_ID", "HuggingFaceTB/SmolLM2-135M-Instruct")


@dataclass
class ChatMessage:
    role: str
    content: str


@dataclass
class GenerateRequest:
    messages: list[ChatMessage]
    max_tokens: int = MAX_NEW_TOKENS
    temperature: float = DEFAULT_TEMPERATURE
    model: str = ""
    stream: bool = False


@dataclass
class EngineInfo:
    backend: str
    model_id: str
    device: str
    ready: bool
    load_ms: int = 0
    extra: dict = field(default_factory=dict)


class InferenceEngine(ABC):
    @abstractmethod
    def info(self) -> EngineInfo: ...

    @abstractmethod
    def generate(self, req: GenerateRequest) -> str: ...

    @abstractmethod
    def stream(self, req: GenerateRequest) -> Generator[str, None, None]: ...


class TransformersEngine(InferenceEngine):
    """HuggingFace transformers causal-LM with streaming TextIteratorStreamer."""

    def __init__(self, model_id: str) -> None:
        import torch
        from transformers import AutoModelForCausalLM, AutoTokenizer

        self._model_id = model_id
        self._device = "cuda" if torch.cuda.is_available() else "cpu"
        t0 = time.perf_counter()
        logger.info("loading transformers model %s on %s", model_id, self._device)

        # Local cache first (volume-mounted MODEL_DIR), else Hub download.
        local = MODEL_DIR if (MODEL_DIR / "config.json").exists() else None
        source = str(local) if local else model_id

        self._tok = AutoTokenizer.from_pretrained(source, trust_remote_code=True)
        dtype = torch.float16 if self._device == "cuda" else torch.float32
        self._model = AutoModelForCausalLM.from_pretrained(
            source,
            torch_dtype=dtype,
            trust_remote_code=True,
            low_cpu_mem_usage=True,
        )
        self._model.to(self._device)
        self._model.eval()
        if self._tok.pad_token is None:
            self._tok.pad_token = self._tok.eos_token

        self._load_ms = int((time.perf_counter() - t0) * 1000)
        self._lock = threading.Lock()
        logger.info("model ready in %dms", self._load_ms)

    def info(self) -> EngineInfo:
        return EngineInfo(
            backend="transformers",
            model_id=self._model_id,
            device=self._device,
            ready=True,
            load_ms=self._load_ms,
            extra={"dtype": str(next(self._model.parameters()).dtype)},
        )

    def _prompt(self, messages: list[ChatMessage]) -> str:
        chat = [{"role": m.role, "content": m.content} for m in messages]
        if hasattr(self._tok, "apply_chat_template") and self._tok.chat_template:
            return self._tok.apply_chat_template(
                chat, tokenize=False, add_generation_prompt=True
            )
        # Fallback plain concat for base models without a chat template.
        parts = []
        for m in messages:
            parts.append(f"{m.role}: {m.content}")
        parts.append("assistant:")
        return "\n".join(parts)

    def _gen_kwargs(self, req: GenerateRequest) -> dict:
        import torch

        prompt = self._prompt(req.messages)
        inputs = self._tok(prompt, return_tensors="pt")
        inputs = {k: v.to(self._device) for k, v in inputs.items()}
        max_new = max(1, min(req.max_tokens or MAX_NEW_TOKENS, 2048))
        do_sample = req.temperature is not None and req.temperature > 0
        kwargs = {
            **inputs,
            "max_new_tokens": max_new,
            "do_sample": do_sample,
            "pad_token_id": self._tok.pad_token_id,
            "eos_token_id": self._tok.eos_token_id,
        }
        if do_sample:
            kwargs["temperature"] = max(0.01, float(req.temperature))
            kwargs["top_p"] = 0.9
        return kwargs, inputs["input_ids"].shape[-1]

    def generate(self, req: GenerateRequest) -> str:
        import torch

        with self._lock:
            kwargs, prompt_len = self._gen_kwargs(req)
            with torch.inference_mode():
                out = self._model.generate(**kwargs)
            new_tokens = out[0][prompt_len:]
            return self._tok.decode(new_tokens, skip_special_tokens=True).strip()

    def stream(self, req: GenerateRequest) -> Generator[str, None, None]:
        import torch
        from transformers import TextIteratorStreamer

        # Serialize concurrent generations; streamer still yields tokens ASAP.
        self._lock.acquire()
        try:
            kwargs, _ = self._gen_kwargs(req)
            streamer = TextIteratorStreamer(
                self._tok, skip_prompt=True, skip_special_tokens=True
            )
            kwargs["streamer"] = streamer

            def _run() -> None:
                try:
                    with torch.inference_mode():
                        self._model.generate(**kwargs)
                except Exception:
                    logger.exception("generate failed")
                finally:
                    self._lock.release()

            t = threading.Thread(target=_run, daemon=True)
            t.start()
            for piece in streamer:
                if piece:
                    yield piece
            t.join(timeout=300)
        except Exception:
            self._lock.release()
            raise


class ONNXEngine(InferenceEngine):
    """ONNX Runtime causal generation when a converted model.onnx is mounted.

    Expects either:
      MODEL_DIR/model.onnx + tokenizer files, or
      MODEL_DIR/onnx/model.onnx

    Uses greedy / multinomial decode loop with the HF tokenizer. This is a
    lightweight path for exported SmolLM/ONNX models — not full ORT-GenAI.
    """

    def __init__(self, model_id: str, onnx_path: Path) -> None:
        import onnxruntime as ort
        from transformers import AutoTokenizer

        self._model_id = model_id
        self._onnx_path = onnx_path
        t0 = time.perf_counter()
        logger.info("loading ONNX model from %s", onnx_path)

        tok_src = str(MODEL_DIR) if (MODEL_DIR / "tokenizer_config.json").exists() else model_id
        self._tok = AutoTokenizer.from_pretrained(tok_src, trust_remote_code=True)
        if self._tok.pad_token is None:
            self._tok.pad_token = self._tok.eos_token

        opts = ort.SessionOptions()
        opts.graph_optimization_level = ort.GraphOptimizationLevel.ORT_ENABLE_ALL
        opts.intra_op_num_threads = int(os.getenv("ORT_INTRA_THREADS", "4"))
        providers = ["CPUExecutionProvider"]
        # CUDA EP when the image was built with GPU support.
        avail = ort.get_available_providers()
        if "CUDAExecutionProvider" in avail:
            providers.insert(0, "CUDAExecutionProvider")
        self._sess = ort.InferenceSession(str(onnx_path), opts, providers=providers)
        self._input_names = [i.name for i in self._sess.get_inputs()]
        self._output_names = [o.name for o in self._sess.get_outputs()]
        self._device = self._sess.get_providers()[0]
        self._load_ms = int((time.perf_counter() - t0) * 1000)
        self._lock = threading.Lock()
        logger.info(
            "ONNX ready in %dms providers=%s inputs=%s",
            self._load_ms,
            self._sess.get_providers(),
            self._input_names,
        )

    def info(self) -> EngineInfo:
        return EngineInfo(
            backend="onnxruntime",
            model_id=self._model_id,
            device=self._device,
            ready=True,
            load_ms=self._load_ms,
            extra={
                "onnx_path": str(self._onnx_path),
                "providers": self._sess.get_providers(),
            },
        )

    def _prompt_ids(self, messages: list[ChatMessage]):
        import numpy as np

        chat = [{"role": m.role, "content": m.content} for m in messages]
        if hasattr(self._tok, "apply_chat_template") and self._tok.chat_template:
            text = self._tok.apply_chat_template(
                chat, tokenize=False, add_generation_prompt=True
            )
        else:
            text = "\n".join(f"{m.role}: {m.content}" for m in messages) + "\nassistant:"
        ids = self._tok.encode(text, return_tensors=None)
        return np.array([ids], dtype=np.int64)

    def _forward_logits(self, input_ids):
        import numpy as np

        feeds = {}
        # Common export shapes: input_ids, attention_mask, position_ids
        if "input_ids" in self._input_names:
            feeds["input_ids"] = input_ids
        else:
            feeds[self._input_names[0]] = input_ids
        if "attention_mask" in self._input_names:
            feeds["attention_mask"] = np.ones_like(input_ids)
        if "position_ids" in self._input_names:
            feeds["position_ids"] = np.arange(input_ids.shape[1], dtype=np.int64)[None, :]

        outs = self._sess.run(None, feeds)
        logits = outs[0]
        # logits: [batch, seq, vocab] or [batch, vocab]
        if logits.ndim == 3:
            return logits[0, -1]
        return logits[0]

    def _sample(self, logits, temperature: float) -> int:
        import numpy as np

        if temperature is None or temperature <= 0:
            return int(np.argmax(logits))
        logits = logits.astype(np.float64) / max(0.01, temperature)
        logits = logits - logits.max()
        probs = np.exp(logits)
        probs = probs / probs.sum()
        return int(np.random.choice(len(probs), p=probs))

    def generate(self, req: GenerateRequest) -> str:
        return "".join(self.stream(req))

    def stream(self, req: GenerateRequest) -> Generator[str, None, None]:
        import numpy as np

        with self._lock:
            ids = self._prompt_ids(req.messages)
            max_new = max(1, min(req.max_tokens or MAX_NEW_TOKENS, 512))
            eos = self._tok.eos_token_id
            for _ in range(max_new):
                logits = self._forward_logits(ids)
                next_id = self._sample(logits, req.temperature)
                if eos is not None and next_id == eos:
                    break
                ids = np.concatenate([ids, np.array([[next_id]], dtype=np.int64)], axis=1)
                piece = self._tok.decode([next_id], skip_special_tokens=True)
                if piece:
                    yield piece


def _find_onnx(root: Path) -> Optional[Path]:
    candidates = [
        root / "model.onnx",
        root / "model_q4.onnx",
        root / "onnx" / "model.onnx",
        root / "onnx" / "model_q4f16.onnx",
    ]
    for c in candidates:
        if c.is_file():
            return c
    # Any single .onnx under root (depth 2)
    if root.is_dir():
        found = list(root.glob("*.onnx")) + list(root.glob("onnx/*.onnx"))
        if found:
            return found[0]
    return None


_engine: Optional[InferenceEngine] = None
_engine_error: Optional[str] = None


def get_engine() -> InferenceEngine:
    global _engine, _engine_error
    if _engine is not None:
        return _engine
    if _engine_error:
        raise RuntimeError(_engine_error)

    model_id = DEFAULT_HF_ID
    prefer = os.getenv("SLM_BACKEND", "auto").lower()  # auto | transformers | onnx
    onnx_path = _find_onnx(MODEL_DIR)

    try:
        if prefer == "onnx" or (prefer == "auto" and onnx_path is not None):
            if onnx_path is None:
                raise RuntimeError("SLM_BACKEND=onnx but no model.onnx under MODEL_DIR")
            _engine = ONNXEngine(model_id, onnx_path)
            return _engine
        try:
            _engine = TransformersEngine(model_id)
        except Exception as primary_err:
            logger.warning("primary model %s failed: %s — trying fallback", model_id, primary_err)
            if model_id != FALLBACK_HF_ID:
                _engine = TransformersEngine(FALLBACK_HF_ID)
            else:
                raise
        return _engine
    except Exception as e:
        _engine_error = str(e)
        logger.exception("failed to load SLM engine")
        raise


def engine_ready() -> bool:
    try:
        get_engine()
        return True
    except Exception:
        return False


def list_models() -> list[dict]:
    """OpenAI-shaped model list for /v1/models."""
    try:
        eng = get_engine()
        info = eng.info()
        mid = info.model_id
    except Exception:
        mid = DEFAULT_HF_ID
    aliases = {
        mid,
        "smollm",
        "smollm2",
        "smollm2-360m",
        "smollm2-135m",
        Path(mid).name if "/" in mid else mid,
    }
    return [{"id": m, "object": "model", "owned_by": "talon-slm"} for m in sorted(aliases)]
