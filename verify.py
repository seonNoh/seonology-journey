#!/usr/bin/env python3
"""Verify the GTM-LIVE-65 repository contract without exposing credentials."""

from __future__ import annotations

import argparse
import re
import subprocess
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parent
BASELINE = "c076de3a6f1e3e298bea47f0474935389203d6b3"
EMOJI = re.compile("[\U0001F1E6-\U0001F1FF\U0001F300-\U0001FAFF\U00002600-\U000027BF]")
AI_SIGNATURE = re.compile(
    r"(?i)(co-authored-by:.*(?:codex|openai|anthropic|claude)|generated (?:with|by) (?:codex|openai|anthropic|claude))"
)


def run(command: list[str], cwd: Path = ROOT) -> bool:
    return subprocess.run(command, cwd=cwd, check=False).returncode == 0


def migration_files() -> list[Path]:
    paths = [
        ROOT / "README.md",
        ROOT / "README.ko.md",
        ROOT / "README.ja.md",
        ROOT / "README_STRUCTURE.md",
        ROOT / "CONTRIBUTING.md",
    ]
    for directory in (ROOT / ".gitea", ROOT / "docs" / "svg"):
        paths.extend(path for path in directory.rglob("*") if path.is_file())
    return paths


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--skip-components", action="store_true")
    args = parser.parse_args()
    failures: list[str] = []

    def check(name: str, passed: bool) -> None:
        print(f"[{'PASS' if passed else 'FAIL'}] {name}")
        if not passed:
            failures.append(name)

    check("repository contract", run([sys.executable, "-m", "unittest", "tests/test_repository_contract.py"]))
    check("deterministic diagrams", run([sys.executable, "tools/generate_phase_b_svgs.py", "--check"]))
    check("GitHub workflow baseline", run(["git", "diff", "--exit-code", BASELINE, "--", ".github/workflows"]))

    emoji_hits = []
    signature_hits = []
    for path in migration_files():
        body = path.read_text(encoding="utf-8", errors="ignore")
        if EMOJI.search(body):
            emoji_hits.append(str(path.relative_to(ROOT)))
        if AI_SIGNATURE.search(body):
            signature_hits.append(str(path.relative_to(ROOT)))
    check("emoji policy", not emoji_hits)
    check("automated-author signature policy", not signature_hits)

    if not args.skip_components:
        for component in ("api", "back"):
            cwd = ROOT / "apps" / component
            check(f"{component} vet", run(["go", "vet", "./..."], cwd))
            check(f"{component} test", run(["go", "test", "./...", "-race", "-count=1"], cwd))
            check(f"{component} build", run(["go", "build", "./..."], cwd))
        check("web typecheck", run(["pnpm", "--filter", "@seonology/journey-web", "typecheck"]))
        check("web test", run(["pnpm", "--filter", "@seonology/journey-web", "test"]))
        check("web build", run(["pnpm", "--filter", "@seonology/journey-web", "build"]))

    print("SUMMARY:", "ALL PASS" if not failures else f"{len(failures)} FAIL")
    return 1 if failures else 0


if __name__ == "__main__":
    raise SystemExit(main())
