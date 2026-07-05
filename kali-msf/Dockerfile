# Use the official lightweight Kali Rolling image
FROM kalilinux/kali-rolling

# Update repositories and install metasploit-framework
RUN apt-get update && \
    apt-get install -y --no-install-recommends metasploit-framework && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Set a default port in case the environment variable is missing
ENV MSF_PORT=5554

# Run msfrpcd in the foreground using the port from the environment variable
# 'sh -c' is required here so Docker evaluates the $MSF_PORT variable at runtime
CMD sh -c "msfrpcd -P network_msf -f -a 0.0.0.0 -p $MSF_PORT"
