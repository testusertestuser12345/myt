
pip install mythic


#!/usr/bin/env python3
"""
Fetch Mythic command history for a given callback and print it as
readable, formatted text.

Usage:
    python3 history.py <callback_number> [-o output_file]
Examples:
    python3 history.py 1              # print to screen
    python3 history.py 1 -o cb1.txt   # save to file

Run "pip install mythic" first.
"""
import argparse
import asyncio
import base64
import binascii
from mythic import mythic

# ==== edit for your environment ====
USERNAME    = "mythic_admin"
PASSWORD    = "put MYTHIC_ADMIN_PASSWORD from Mythic/.env here"
SERVER_IP   = "127.0.0.1"
SERVER_PORT = 7443          # Mythic default port
PAGE        = 500           # rows fetched per request
# ===================================

QUERY = """
query History($cb: Int!, $limit: Int!, $offset: Int!) {
  task(
    where: {callback: {display_id: {_eq: $cb}}},
    order_by: {id: asc},
    limit: $limit, offset: $offset
  ) {
    display_id
    command_name
    display_params
    status
    timestamp
    operator { username }
    responses(order_by: {id: asc}) { response_text }
  }
}
"""

LINE = "-" * 60


def decode_chunk(s):
    """Decode base64 output. Return as-is if it is not base64."""
    if not s:
        return ""
    t = "".join(s.split())          # strip all whitespace/newlines
    t += "=" * (-len(t) % 4)        # fix padding
    try:
        return base64.b64decode(t).decode("utf-8", errors="replace")
    except (binascii.Error, ValueError):
        return s


def format_task(t, seq):
    """Format one task as a readable text block (seq = per-callback index)."""
    cmd = t.get("command_name", "")
    params = t.get("display_params", "") or ""
    full_cmd = f"{cmd} {params}".strip()

    # join response chunks, base64-decoded
    chunks = [decode_chunk(r.get("response_text")) for r in (t.get("responses") or [])]
    result = "".join(chunks).rstrip()
    if not result:
        result = "(no output)"

    head = (f"#{seq}  (id:{t.get('display_id','')})  |  {t.get('timestamp','')}  |  "
            f"{t.get('status','')}  |  {(t.get('operator') or {}).get('username','')}")

    return (
        f"{LINE}\n"
        f"{head}\n"
        f"cmd: {full_cmd}\n"
        f"{LINE}\n"
        f"{result}\n"
    )


async def fetch_all(inst, cb):
    rows = []
    offset = 0
    while True:
        data = await mythic.execute_custom_query(
            inst, query=QUERY,
            variables={"cb": cb, "limit": PAGE, "offset": offset},
        )
        batch = data["task"]
        if not batch:
            break
        rows.extend(batch)
        offset += PAGE
    return rows


async def main():
    parser = argparse.ArgumentParser(description="Print Mythic command history as text")
    parser.add_argument("callback", type=int, help="callback number")
    parser.add_argument("-o", "--output", help="output file (prints to screen if omitted)")
    args = parser.parse_args()

    inst = await mythic.login(
        username=USERNAME,
        password=PASSWORD,
        server_ip=SERVER_IP,
        server_port=SERVER_PORT,
        ssl=True,
        timeout=-1,
    )

    tasks = await fetch_all(inst, args.callback)

    title = f"=== Callback {args.callback} command history ({len(tasks)} tasks) ===\n\n"
    body = "\n".join(format_task(t, i) for i, t in enumerate(tasks, 1))
    text = title + body

    if args.output:
        with open(args.output, "w", encoding="utf-8") as f:
            f.write(text)
        print(f"Done: wrote {len(tasks)} tasks to {args.output}")
    else:
        print(text)


if __name__ == "__main__":
    asyncio.run(main())
