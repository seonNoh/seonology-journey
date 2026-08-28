from __future__ import annotations

import hashlib
import json
import re
import subprocess
import sys
import unittest
import xml.etree.ElementTree as ET
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
BASELINE = "c076de3a6f1e3e298bea47f0474935389203d6b3"
TOPICS = ("service-architecture", "component-delivery", "runtime-cutover", "repository-roles")
LANGUAGES = ("en", "ko", "ja")
SVG_FILES = {
    f"{topic}{'' if language == 'en' else f'.{language}'}.svg"
    for topic in TOPICS
    for language in LANGUAGES
}
WORKFLOW_HASHES = {
    ".github/workflows/android-ci.yml": "67a95865d4fa1691e5d985d28e4d962d92026ee1acca2a6308bc1c7adf0eb9c6",
    ".github/workflows/api-ci.yml": (
        "6dda1766d82dd2161e29570d6f732601"
        "11247a58df8eea87b41e283b8dfd1114"
    ),
    ".github/workflows/back-ci.yml": "bb9bd53e555079ae36995441edb82ac5da15e0614b2517aec3130825f4838941",
    ".github/workflows/e2e-web.yml": "6a95ad66f4be7e8b6296658209eeeea931e916b2a8550bf55ff0b889de1461e2",
    ".github/workflows/image-scan.yml": "f637d6f7394e68dad4b3a2716a4fe1a4f65bb77caeb5732d55ddb95ec05a40a9",
    ".github/workflows/proto-check.yml": "29efd6bfa203a1c56cad86cbd141dfda41aa0f8a1855e14d920cd95600087b9c",
    ".github/workflows/release.yml": "af14164e40ad3963ba0aa0ea1ca948a7d0b422e2ae53dbe3b19bb1a627b007ba",
    ".github/workflows/security.yml": "f48df44eed71a704ba5c1cb10c98f219720fdbf2358d0970563a49c3f8657fbe",
    ".github/workflows/web-ci.yml": "bb640f9dc078e06b2b55801c0d0b79284ca6074fb210c5d314a61b5fff72fb8d",
    ".github/workflows/zap-scan.yml": "88dac4199cbf25f352d9a731f912ce521cb3ad1a87b431ed7b227411059d08c9",
}


class RepositoryContractTest(unittest.TestCase):
    def test_required_phase_b_files_exist(self) -> None:
        required = {
            "README.md",
            "README.ko.md",
            "README.ja.md",
            "README_STRUCTURE.md",
            "CONTRIBUTING.md",
            "LICENSE",
            ".editorconfig",
            ".gitea/ISSUE_TEMPLATE/bug-report.yaml",
            ".gitea/ISSUE_TEMPLATE/feature-request.yaml",
            ".gitea/PULL_REQUEST_TEMPLATE.md",
            ".gitea/workflows/ci.yml",
            ".gitea/workflows/images.yml",
            ".gitea/workflows/gitea-android-ci.yml",
            ".gitea/workflows/gitea-release.yml",
            "assets/diagrams/pedia-manifest.json",
            "tools/generate_phase_b_svgs.py",
            "verify.py",
        }
        self.assertEqual([], sorted(path for path in required if not (ROOT / path).is_file()))

    def test_readmes_have_equal_structure_and_diagrams(self) -> None:
        names = ("README.md", "README.ko.md", "README.ja.md")
        switcher = "[English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md)"
        heading_shapes = []
        diagrams = []
        for name in names:
            body = (ROOT / name).read_text(encoding="utf-8")
            self.assertIn(switcher, body.splitlines()[:5])
            heading_shapes.append(tuple(len(match.group(1)) for match in re.finditer(r"^(#+) ", body, re.MULTILINE)))
            diagrams.append(
                tuple(
                    re.sub(r"\.(?:ko|ja)\.svg$", ".svg", item)
                    for item in re.findall(r"!\[[^]]*]\((assets/diagrams/[^)]+)\)", body)
                )
            )
        self.assertEqual(heading_shapes[0], heading_shapes[1])
        self.assertEqual(heading_shapes[0], heading_shapes[2])
        self.assertEqual(diagrams[0], diagrams[1])
        self.assertEqual(diagrams[0], diagrams[2])
        self.assertEqual(4, len(diagrams[0]))

    def test_relief_svg_set_is_self_contained_and_deterministic(self) -> None:
        svg_dir = ROOT / "assets" / "diagrams"
        actual = {path.name for path in svg_dir.glob("*.svg")}
        self.assertEqual(SVG_FILES, actual)
        all_ids: set[str] = set()
        forbidden = re.compile(r"#(?:0f172a|1e293b|38bdf8|a78bfa|f472b6|34d399|fbbf24)", re.IGNORECASE)
        for name in sorted(actual):
            body = (svg_dir / name).read_text(encoding="utf-8")
            ET.fromstring(body)
            for token in ("<style>", "<defs>", "prefers-reduced-motion", "#0d1117", "#7c9fff"):
                self.assertIn(token, body, name)
            self.assertIsNone(forbidden.search(body), name)
            self.assertNotRegex(body, r"<animate(?:Motion)?\b", name)
            ids = re.findall(r'\bid="([^"]+)"', body)
            self.assertEqual(len(ids), len(set(ids)), name)
            self.assertTrue(all_ids.isdisjoint(ids), name)
            all_ids.update(ids)
        result = subprocess.run(
            [sys.executable, "tools/generate_phase_b_svgs.py", "--check"],
            cwd=ROOT,
            text=True,
            capture_output=True,
            check=False,
        )
        self.assertEqual(0, result.returncode, result.stdout + result.stderr)

    def test_pedia_manifest_has_twelve_stable_ids(self) -> None:
        manifest = json.loads(
            (ROOT / "assets/diagrams/pedia-manifest.json").read_text(encoding="utf-8")
        )
        entries = manifest["entries"]
        identifiers = [entry["pedia_id"] for entry in entries]
        self.assertEqual(12, len(entries))
        self.assertEqual(12, len(set(identifiers)))
        self.assertEqual(SVG_FILES, {entry["file"] for entry in entries})
        self.assertTrue(
            all(value.startswith("svg-svg-gtm-live-65-journey-") for value in identifiers)
        )

    def test_github_workflow_bytes_are_preserved(self) -> None:
        for path, expected in WORKFLOW_HASHES.items():
            self.assertEqual(expected, hashlib.sha256((ROOT / path).read_bytes()).hexdigest(), path)
        result = subprocess.run(
            ["git", "diff", "--exit-code", BASELINE, "--", ".github/workflows"],
            cwd=ROOT,
            check=False,
        )
        self.assertEqual(0, result.returncode)

    def test_gitea_workflows_are_native_and_runner_safe(self) -> None:
        workflow_dir = ROOT / ".gitea" / "workflows"
        bodies = {path.name: path.read_text(encoding="utf-8") for path in workflow_dir.glob("*.yml")}
        joined = "\n".join(bodies.values())
        for token in (
            "https://gitea.com/actions/checkout@v4",
            "git.seonology.com/seonology/seonology-journey-api",
            "git.seonology.com/seonology/seonology-journey-back",
            "git.seonology.com/seonology/seonology-journey-web",
            "linux/amd64,linux/arm64",
            "secrets.PACKAGE_USER",
            "secrets.PACKAGE_TOKEN",
        ):
            self.assertIn(token, joined)
        for forbidden in ("ghcr.io/seonnoh", "GITHUB_TOKEN", "api.github.com", "actions/setup-go", "actions/setup-node"):
            self.assertNotIn(forbidden, joined)
        self.assertNotIn("matrix:", bodies["images.yml"])
        self.assertIn("for spec in", bodies["images.yml"])
        self.assertEqual(3, bodies["images.yml"].count("git.seonology.com/seonology/seonology-journey-"))
        self.assertIn("fetch-depth: 0", bodies["ci.yml"])

    def test_policy_for_new_migration_files(self) -> None:
        emoji = re.compile("[\U0001F1E6-\U0001F1FF\U0001F300-\U0001FAFF\U00002600-\U000027BF]")
        signature = re.compile(
            r"(?i)(co-authored-by:.*(?:codex|openai|anthropic|claude)|generated (?:with|by) (?:codex|openai|anthropic|claude))"
        )
        targets = [
            ROOT / "README.md",
            ROOT / "README.ko.md",
            ROOT / "README.ja.md",
            ROOT / "README_STRUCTURE.md",
            ROOT / "CONTRIBUTING.md",
        ]
        targets.extend((ROOT / ".gitea").rglob("*"))
        targets.extend((ROOT / "assets" / "diagrams").glob("*.svg"))
        for path in (item for item in targets if item.is_file()):
            body = path.read_text(encoding="utf-8", errors="ignore")
            self.assertIsNone(emoji.search(body), str(path.relative_to(ROOT)))
            self.assertIsNone(signature.search(body), str(path.relative_to(ROOT)))

    def test_vitest_excludes_playwright_e2e_specs(self) -> None:
        body = (ROOT / "apps/web/vite.config.ts").read_text(encoding="utf-8")
        self.assertIn("exclude:", body)
        self.assertIn("e2e/**", body)

    def test_docker_context_excludes_local_build_artifacts(self) -> None:
        self.assertTrue((ROOT / ".dockerignore").is_file())
        entries = {
            line.strip()
            for line in (ROOT / ".dockerignore").read_text(encoding="utf-8").splitlines()
            if line.strip() and not line.startswith("#")
        }
        self.assertTrue({".git", "**/node_modules", "**/dist", "**/build"}.issubset(entries))


if __name__ == "__main__":
    unittest.main()
