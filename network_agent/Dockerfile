# Use the specific Python version you requested
FROM python:3.13-slim

# Set the working directory inside the container
WORKDIR /app

# Install system dependencies (optional but recommended for network tools)
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    python3-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements first to leverage Docker cache
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of your application code
COPY . .

# Keep the container running for interactive use or script execution
# Keep the container running for interactive use or script execution
 CMD ["uvicorn", "pentest_core.fast_api:app", "--host", "0.0.0.0", "--port", "8000"]
