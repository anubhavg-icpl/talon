# SmolLM / ONNX model mount

This directory is volume-mounted into the `onnx-slm` container at `/models/smollm`.

## Default (transformers / Hub)

No files required. On first start the service downloads:

- `HuggingFaceTB/SmolLM2-360M-Instruct` (default)
- or `HuggingFaceTB/SmolLM2-135M-Instruct` (set `SLM_MODEL_ID` or fallback)

Weights land in `./models/hf-cache` (compose volume) so restarts are fast.

## Optional ONNX Runtime path

Place an exported ONNX graph here and set `SLM_BACKEND=onnx` (or `auto`):

```
models/smollm/
  model.onnx          # or model_q4.onnx / onnx/model.onnx
  tokenizer files…    # tokenizer.json / tokenizer_config.json / …
```

The service uses `onnxruntime` with a greedy/multinomial decode loop.

## Ollama alternative

For GGUF via Ollama instead of this runtime:

```bash
docker compose --profile ollama up -d
ollama pull smollm2:360m
# LLM_PROVIDER=ollama  OLLAMA_MAIN_MODEL=smollm2:360m
```
