#!/usr/bin/env python3
"""Generate the twelve deterministic Relief diagrams for GTM-LIVE-65."""

from __future__ import annotations

import argparse
import html
import json
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
OUT = ROOT / "assets" / "diagrams"
LANG_SUFFIX = {"en": "", "ko": ".ko", "ja": ".ja"}

DATA = {
    "service-architecture": {
        "en": ("Service architecture", "Web and Android clients share two Go services", ("Clients", "Web and Android", "Journey API", "REST and WebSocket", "Journey back", "gRPC and AWS data")),
        "ko": ("서비스 아키텍처", "웹과 Android 클라이언트가 두 Go 서비스를 공유합니다", ("클라이언트", "웹과 Android", "Journey API", "REST와 WebSocket", "Journey back", "gRPC와 AWS 데이터")),
        "ja": ("サービス構成", "WebとAndroidが2つのGoサービスを共有します", ("クライアント", "WebとAndroid", "Journey API", "RESTとWebSocket", "Journey back", "gRPCとAWSデータ")),
    },
    "component-delivery": {
        "en": ("Component delivery", "One runner publishes three images serially", ("Validated source", "One main commit", "Gitea Actions", "max-parallel one", "OCI Registry", "api, back and web")),
        "ko": ("컴포넌트 배포", "단일 러너가 이미지 3개를 직렬로 게시합니다", ("검증된 소스", "단일 main 커밋", "Gitea Actions", "동시 실행 한 개", "OCI Registry", "api, back, web")),
        "ja": ("コンポーネント配布", "単一ランナーが3イメージを順番に公開します", ("検証済みソース", "単一mainコミット", "Gitea Actions", "同時実行は1件", "OCI Registry", "api、back、web")),
    },
    "runtime-cutover": {
        "en": ("Runtime cutover", "GitOps pins Gitea images before health approval", ("Gitea Registry", "Three multi-arch images", "Argo CD", "Synced and Healthy", "k3s runtime", "Six pods ready")),
        "ko": ("운영 전환", "상태 승인 전에 GitOps가 Gitea 이미지를 고정합니다", ("Gitea Registry", "다중 아키텍처 이미지 3개", "Argo CD", "Synced와 Healthy", "k3s 런타임", "Pod 6개 Ready")),
        "ja": ("本番切替", "状態承認前にGitOpsでGiteaイメージを固定します", ("Gitea Registry", "マルチアーキ3イメージ", "Argo CD", "SyncedとHealthy", "k3s実行環境", "6 PodがReady")),
    },
    "repository-roles": {
        "en": ("Repository roles", "Gitea is the source of truth and GitHub is a mirror", ("Contributors", "Ordinary main changes", "Gitea", "Source of truth", "GitHub", "Read-only push mirror")),
        "ko": ("저장소 역할", "Gitea는 기준 저장소이고 GitHub는 미러입니다", ("기여자", "일반 main 변경", "Gitea", "기준 저장소", "GitHub", "읽기 전용 push mirror")),
        "ja": ("リポジトリ役割", "Giteaを正本としGitHubをミラーにします", ("コントリビューター", "通常のmain変更", "Gitea", "信頼できる正本", "GitHub", "読取専用push mirror")),
    },
}


def render(topic: str, language: str) -> str:
    title, subtitle, labels = DATA[topic][language]
    prefix = f"g65-{topic}-{language}"
    cards = ((56, labels[0], labels[1]), (352, labels[2], labels[3]), (648, labels[4], labels[5]))
    markup = []
    for index, (x, heading, detail) in enumerate(cards):
        lead = index == 1
        markup.append(
            f'''  <g filter="url(#{prefix}-{'lead-shadow' if lead else 'shadow'})">
    <rect x="{x}" y="144" width="256" height="224" rx="12" fill="url(#{prefix}-{'lead-surface' if lead else 'surface'})" stroke="{'#33404f' if lead else '#252d38'}"/>
    <path d="M{x + 12} 145 H{x + 244}" stroke="{'#7c9fff' if lead else '#ffffff'}" stroke-opacity="{'0.40' if lead else '0.05'}"/>
    <rect x="{x + 18}" y="170" width="220" height="38" rx="8" fill="{'#232c39' if lead else '#1e2530'}"/>
    <text x="{x + 128}" y="194" class="card" text-anchor="middle">{html.escape(heading)}</text>
    <text x="{x + 128}" y="246" class="body" text-anchor="middle">{html.escape(detail)}</text>
    <rect x="{x + 42}" y="286" width="172" height="30" rx="8" fill="#151b23" stroke="#252d38"/>
    <text x="{x + 128}" y="306" class="value" text-anchor="middle">{'source' if index == 0 else 'verified' if index == 1 else 'ready'}</text>
  </g>'''
        )
    return f'''<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 960 540" role="img" aria-labelledby="{prefix}-title {prefix}-desc">
  <title id="{prefix}-title">{html.escape(title)}</title>
  <desc id="{prefix}-desc">{html.escape(subtitle)}</desc>
  <style>
    .title {{ fill:#f0f6fc; font:600 24px "Pretendard",system-ui,sans-serif; }}
    .subtitle {{ fill:#939da8; font:400 13px "Pretendard",system-ui,sans-serif; }}
    .card {{ fill:#f0f6fc; font:500 15px "Pretendard",system-ui,sans-serif; }}
    .body {{ fill:#c3ccd6; font:400 12px "Pretendard",system-ui,sans-serif; }}
    .value {{ fill:#7c9fff; font:500 11px ui-monospace,monospace; }}
    .foot {{ fill:#69717a; font:400 11px "Pretendard",system-ui,sans-serif; }}
    .flow {{ fill:none; stroke:#7c9fff; stroke-width:1.6; stroke-dasharray:5 6; animation:{prefix}-flow 2s linear infinite; }}
    @keyframes {prefix}-flow {{ to {{ stroke-dashoffset:-22; }} }}
    @media (prefers-reduced-motion: reduce) {{ .flow {{ animation:none; }} }}
  </style>
  <defs>
    <linearGradient id="{prefix}-surface" x1="0" y1="0" x2="0" y2="1"><stop stop-color="#1b222c"/><stop offset="1" stop-color="#151b23"/></linearGradient>
    <linearGradient id="{prefix}-lead-surface" x1="0" y1="0" x2="0" y2="1"><stop stop-color="#222b38"/><stop offset="1" stop-color="#19202a"/></linearGradient>
    <filter id="{prefix}-shadow" x="-20%" y="-20%" width="150%" height="170%"><feDropShadow dx="-3" dy="10" stdDeviation="9" flood-color="#000000" flood-opacity=".50"/></filter>
    <filter id="{prefix}-lead-shadow" x="-20%" y="-20%" width="150%" height="180%"><feDropShadow dx="-4" dy="13" stdDeviation="11" flood-color="#000000" flood-opacity=".58"/></filter>
    <marker id="{prefix}-arrow" viewBox="0 0 8 8" refX="7" refY="4" markerWidth="6" markerHeight="6" orient="auto"><path d="M0,0 L8,4 L0,8 z" fill="#7c9fff"/></marker>
  </defs>
  <rect width="960" height="540" fill="#0d1117"/>
  <text x="48" y="56" class="title">{html.escape(title)}</text>
  <text x="48" y="82" class="subtitle">{html.escape(subtitle)}</text>
{chr(10).join(markup)}
  <path class="flow" marker-end="url(#{prefix}-arrow)" d="M312 256 H340 Q344 256 344 260 V260"/>
  <path class="flow" marker-end="url(#{prefix}-arrow)" d="M608 256 H636 Q640 256 640 260 V260"/>
  <text x="48" y="474" class="foot">GTM-LIVE-65 · seonology journey · Gitea SSOT</text>
</svg>
'''


def outputs() -> dict[Path, str]:
    result: dict[Path, str] = {}
    entries = []
    for topic in DATA:
        for language in LANG_SUFFIX:
            filename = f"{topic}{LANG_SUFFIX[language]}.svg"
            result[OUT / filename] = render(topic, language)
            entries.append(
                {
                    "file": filename,
                    "pedia_id": f"svg-svg-gtm-live-65-journey-{topic}-{language}",
                }
            )
    result[OUT / "pedia-manifest.json"] = json.dumps({"contract_version": 1, "entries": entries}, ensure_ascii=False, indent=2) + "\n"
    return result


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true")
    args = parser.parse_args()
    mismatches = []
    for path, expected in outputs().items():
        if args.check:
            actual = path.read_text(encoding="utf-8") if path.is_file() else None
            if actual != expected:
                mismatches.append(str(path.relative_to(ROOT)))
        else:
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_text(expected, encoding="utf-8")
    if mismatches:
        print("Generated files differ:", ", ".join(mismatches))
        return 1
    print("Phase B SVG set is deterministic" if args.check else "Generated Phase B SVG set")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
