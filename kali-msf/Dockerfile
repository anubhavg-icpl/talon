# Use the official lightweight Kali Rolling image
FROM kalilinux/kali-rolling

# Update repositories and install metasploit-framework
RUN apt-get update && \
    apt-get install -y --no-install-recommends metasploit-framework && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Set a default port in case the environment variable is missing
ENV MSF_PORT=5554

# MSF_PASSWORD has no default -- msfrpcd refuses to start without one, which
# is the point: no hardcoded credential baked into the image.
# 'sh -c' is required here so Docker evaluates the env vars at runtime
CMD sh -c "if [ -z \"$MSF_PASSWORD\" ]; then echo 'MSF_PASSWORD is required' >&2; exit 1; fi; exec msfrpcd -P \"$MSF_PASSWORD\" -f -a 0.0.0.0 -p $MSF_PORT"
