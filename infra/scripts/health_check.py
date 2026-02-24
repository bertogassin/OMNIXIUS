#!/usr/bin/env python3
"""OMNIXIUS — health check script (Python)."""
import urllib.request
import sys

endpoints = [
    ("Go API", "http://localhost:8080/health"),
    ("Rust", "http://localhost:8081/health"),
    ("AI", "http://localhost:8000/health"),
]

def main():
    ok = 0
    for name, url in endpoints:
        try:
            with urllib.request.urlopen(url, timeout=2) as r:
                if r.status == 200:
                    print(f"OK {name} {url}")
                    ok += 1
                else:
                    print(f"FAIL {name} {url} status={r.status}")
        except Exception as e:
            print(f"FAIL {name} {url} — {e}")
    print(f"{ok}/{len(endpoints)} services up")
    sys.exit(0 if ok == len(endpoints) else 1)

if __name__ == "__main__":
    main()
