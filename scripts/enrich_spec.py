#!/usr/bin/env python3
"""
enrich_spec.py: Inject operationIds into the Bitbucket OpenAPI spec.

Strategy: slugify(summary) if present, else "{method}_{path_slug}".
This is deterministic and stable as long as Atlassian doesn't rename summaries
(rare) — and the CI diff check will catch it if they do.

Usage: python3 enrich_spec.py <input.json> <output.json>
"""

import json
import re
import sys
from pathlib import Path


def to_camel(s: str) -> str:
    """'List pull requests' -> 'listPullRequests'"""
    words = re.sub(r"[^a-zA-Z0-9 ]", "", s).split()
    if not words:
        return ""
    return words[0].lower() + "".join(w.title() for w in words[1:])


def path_slug(path: str, method: str) -> str:
    """/repositories/{workspace}/{repo_slug}/pullrequests' + 'get' -> 'getRepositoriesPullrequests'"""
    parts = [p for p in path.split("/") if p and not p.startswith("{")]
    return method.lower() + "".join(p.title() for p in parts)


def main():
    if len(sys.argv) != 3:
        print(f"Usage: {sys.argv[0]} <input.json> <output.json>", file=sys.stderr)
        sys.exit(1)

    input_path = Path(sys.argv[1])
    output_path = Path(sys.argv[2])

    spec = json.loads(input_path.read_text())

    http_methods = ["get", "post", "put", "patch", "delete"]
    count = 0

    for path, path_item in spec.get("paths", {}).items():
        for method in http_methods:
            op = path_item.get(method)
            if not op:
                continue
            if "operationId" not in op:
                summary = op.get("summary", "")
                op["operationId"] = to_camel(summary) if summary else path_slug(path, method)
            count += 1

    output_path.write_text(json.dumps(spec, indent=2))
    print(f"Enriched {count} operations, wrote to {output_path}")


if __name__ == "__main__":
    main()
