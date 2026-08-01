# Talon SLM runtime — SmolLM2 (+ optional ONNX Runtime) with OpenAI-compatible SSE.
# Profile: docker compose --profile slm up -d --build onnx-slm
#
# CPU image by default. GPU: rebuild with --build-arg TORCH_INDEX=.../cu124 and
# install onnxruntime-gpu instead of onnxruntime.
FROM python:3.12-slim-bookworm

ARG TORCH_INDEX=https://download.pytorch.org/whl/cpu
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    # HF cache inside the volume so restarts skip re-download.
    HF_HOME=/models/hf-cache \
    TRANSFORMERS_CACHE=/models/hf-cache \
    MODEL_DIR=/models/smollm \
    SLM_PORT=8090 \
    SLM_MODEL_ID=HuggingFaceTB/SmolLM2-360M-Instruct \
    SLM_BACKEND=auto \
    SLM_LAZY=0 \
    # Cap concurrent downloads / threads for small hosts.
    OMP_NUM_THREADS=4 \
    ORT_INTRA_THREADS=4

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
        curl ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Torch first (large, cached layer), then the rest.
COPY requirements.txt .
RUN pip install --no-cache-dir torch --index-url ${TORCH_INDEX} \
    && pip install --no-cache-dir -r requirements.txt

COPY app ./app

# Models dir is typically volume-mounted; create stubs so the image runs alone.
RUN mkdir -p /models/smollm /models/hf-cache

EXPOSE 8090

HEALTHCHECK --interval=15s --timeout=5s --start-period=120s --retries=5 \
    CMD curl -fsS http://127.0.0.1:8090/health || exit 1

# Preload model in the app lifespan (SLM_LAZY=0). First start may take minutes
# while HF weights download into the volume.
CMD ["python", "-m", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8090", "--workers", "1"]
