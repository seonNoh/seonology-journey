#!/usr/bin/env python3
"""
PDF 샘플 3건(아칸·아사히카와·오비히로) 을 testuser 계정으로 시드하는 스크립트.

전제:
  - Keycloak: https://auth.seonology.com/realms/seonology
  - API:      https://journey-api.seonology.com
  - testuser / Test1234! 가 존재
  - journey-web 클라이언트(public, PKCE) 의 Direct Access Grants 허용 여부와 무관하게
    journey-android 또는 journey-cli 같은 public 클라이언트에 ROPC 가 필요할 수 있음.
    여기서는 journey-web 의 token 엔드포인트로 password grant 를 시도하고,
    실패 시 사용자가 KEYCLOAK_TOKEN 환경변수로 직접 access_token 을 주입할 수 있다.
"""
from __future__ import annotations

import json
import os
import sys
import time
from typing import Any

import requests

API_BASE = os.environ.get("API_BASE", "https://journey-api.seonology.com")
KC_BASE = os.environ.get(
    "KEYCLOAK_BASE", "https://auth.seonology.com/realms/seonology-journey"
)
USER = os.environ.get("KC_USER", "seon")
PASS = os.environ.get("KC_PASS", "!Tjs78xor0512")
CLIENT_ID = os.environ.get("KC_CLIENT", "journey-web")


def get_token() -> str:
    if t := os.environ.get("KEYCLOAK_TOKEN"):
        return t
    r = requests.post(
        f"{KC_BASE}/protocol/openid-connect/token",
        data={
            "grant_type": "password",
            "client_id": CLIENT_ID,
            "username": USER,
            "password": PASS,
            "scope": "openid",
        },
        timeout=15,
    )
    if r.status_code != 200:
        sys.exit(
            f"keycloak token failed: {r.status_code} {r.text}\n"
            "→ KEYCLOAK_TOKEN 환경변수로 access_token 을 주입하거나, "
            "journey-web 클라이언트의 Direct Access Grants 를 켜라."
        )
    return r.json()["access_token"]


def api(token: str, method: str, path: str, body: dict[str, Any] | None = None) -> Any:
    r = requests.request(
        method,
        f"{API_BASE}{path}",
        headers={
            "Authorization": f"Bearer {token}",
            "Content-Type": "application/json",
        },
        data=json.dumps(body) if body is not None else None,
        timeout=20,
    )
    if r.status_code >= 400:
        raise RuntimeError(f"{method} {path} -> {r.status_code} {r.text}")
    if r.text == "":
        return None
    return r.json()


# === 3개 PDF 데이터 (요약·정형화) ===========================================
# PDF 원본은 source/*.pdf 참조. 본 시드는 핵심 일정/식사/숙박만 적재한다.

TRIPS: list[dict[str, Any]] = [
    {
        "title": "아칸 여행 (마리모 마츠리)",
        "description": "삿포로 출발 → 쿠시로 → 아칸 호숫가 료칸. 마리모 마츠리·유람선·카무이 루미나.",
        "start_date": "2025-10-09",
        "end_date": "2025-10-12",
        "destination": "홋카이도 / 아칸·쿠시로",
        "country_code": "JP",
        "total_budget": {"amount_minor": "350000", "currency": "JPY"},
        "days": [
            {
                "region": "삿포로 → 쿠시로 → 아칸",
                "schedules": [
                    ("08:00", "일어나기", "REST", "WALK", ""),
                    ("09:30", "삿포로 역으로 출발", "TRANSPORT", "TRAIN", "쾌속 에어포트"),
                    ("11:00", "신치토세공항 도착, 짐 붙이기", "TRANSPORT", "FLIGHT", ""),
                    ("13:20", "신치토세공항 출발 ANA4873", "TRANSPORT", "FLIGHT", "ANA4873"),
                    ("14:05", "쿠시로 공항 도착", "TRANSPORT", "FLIGHT", ""),
                    ("15:37", "쿠시로 → 아칸 버스 이동", "TRANSPORT", "BUS", "공항리무진"),
                    ("17:00", "あかん鶴雅別荘 鄙の座 체크인", "ACCOMMODATION", "WALK", ""),
                    ("19:30", "마리모 마츠리 + 온천", "ACTIVITY", "WALK", ""),
                ],
                "meals": [
                    ("LUNCH", "LOCAL", "신치토세공항 식당가"),
                    ("DINNER", "HOTEL", "あかん鶴雅別荘 鄙の座 가이세키"),
                ],
                "accom": ("あかん鶴雅別荘 鄙の座", "17:00", "11:00", "료칸 / 온천 / 가이세키"),
            },
            {
                "region": "아칸 (마리모 마츠리·유람선·카무이 루미나)",
                "schedules": [
                    ("07:30", "호텔 조식 + 온천", "REST", "WALK", ""),
                    ("10:00", "마리모 마츠리 구경", "SIGHTSEEING", "WALK", "10:00 / 10:30 / 11:30"),
                    ("11:00", "阿寒湖畔エコミュージアムセンター", "SIGHTSEEING", "WALK", ""),
                    ("11:30", "ボッケ 遊歩道 (약 1시간)", "ACTIVITY", "WALK", ""),
                    ("13:00", "유람선 탑승 (85분, 1인 2,400엔)", "ACTIVITY", "FERRY", "마리모관람 비용 포함"),
                    ("15:30", "阿寒湖アイヌコタン", "SIGHTSEEING", "WALK", ""),
                    ("16:30", "호텔 복귀 / 온천 휴식", "REST", "WALK", ""),
                    ("19:30", "카무이 루미나 관람 (약 1시간)", "ACTIVITY", "WALK", "전날 호텔에서 티켓 구매"),
                ],
                "meals": [
                    ("BREAKFAST", "HOTEL", "료칸 조식 뷔페"),
                    ("LUNCH", "LOCAL", "아칸호 주변 식당"),
                    ("DINNER", "HOTEL", "료칸 가이세키"),
                ],
                "accom": ("あかん鶴雅別荘 鄙の座", "—", "11:00", "료칸 / 온천 / 가이세키"),
            },
            {
                "region": "아칸 → 쿠시로",
                "schedules": [
                    ("07:30", "호텔 조식 + 온천", "REST", "WALK", ""),
                    ("11:00", "체크아웃 + 온천 마을 산책", "REST", "WALK", ""),
                    ("12:40", "아칸 → 쿠시로 버스 (약 2시간, 2,750엔)", "TRANSPORT", "BUS", "11:50 / 12:40 출발"),
                    ("14:40", "쿠시로 역 도착", "TRANSPORT", "BUS", ""),
                    ("15:00", "쿠시로 스파 호텔 체크인", "ACCOMMODATION", "WALK", ""),
                    ("17:30", "쿠시로 역 근처 산책 + 저녁", "REST", "WALK", ""),
                ],
                "meals": [
                    ("BREAKFAST", "HOTEL", "료칸 조식"),
                    ("LUNCH", "LOCAL", "아칸 온천 마을"),
                    ("DINNER", "LOCAL", "쿠시로 치킨 (난반)"),
                ],
                "accom": ("쿠시로 스파 호텔", "15:00", "10:00", "온천 호텔"),
            },
            {
                "region": "쿠시로 → 삿포로",
                "schedules": [
                    ("09:00", "호텔 조식 후 휴식", "REST", "WALK", ""),
                    ("10:00", "체크아웃", "REST", "WALK", ""),
                    ("12:15", "공항 리무진 버스 (連絡バス12:15, 950엔)", "TRANSPORT", "BUS", "정기 노선 10:10/12:00"),
                    ("13:00", "쿠시로 공항 도착", "TRANSPORT", "BUS", ""),
                    ("14:35", "쿠시로 → 신치토세 ANA4874", "TRANSPORT", "FLIGHT", "ANA4874"),
                    ("15:20", "신치토세공항 도착", "TRANSPORT", "FLIGHT", ""),
                    ("16:00", "쾌속 에어포트 탑승", "TRANSPORT", "TRAIN", ""),
                    ("17:00", "삿포로 역 도착", "TRANSPORT", "TRAIN", ""),
                ],
                "meals": [
                    ("LUNCH", "LOCAL", "쿠시로 부타동"),
                ],
                "accom": None,
            },
        ],
        "expenses": [
            ("ACCOMMODATION", 60000, "JPY", "鶴雅別荘 1박"),
            ("ACCOMMODATION", 60000, "JPY", "鶴雅別荘 2박"),
            ("ACCOMMODATION", 12000, "JPY", "쿠시로 스파 호텔"),
            ("TRANSPORT", 5500, "JPY", "쿠시로↔아칸 버스 왕복"),
            ("ACTIVITY", 4800, "JPY", "유람선 2인"),
            ("FOOD", 8000, "JPY", "외식 합계"),
        ],
        "checklist": [
            ("PACKING", "수영복(온천 X / 가운만 OK)"),
            ("PACKING", "충전기·어댑터"),
            ("BOOKING", "ANA4873 항공 예약"),
            ("BOOKING", "鶴雅別荘 鄙の座 2박"),
            ("BOOKING", "카무이 루미나 티켓"),
            ("TODO", "마리모 마츠리 시간 확인 (10/10)"),
        ],
        "notes": [
            "鶴雅別荘 鄙の座 가이세키는 19:00 시작",
            "유람선은 14:00 표가 마리모 관람 동선상 가장 좋음",
            "카무이 루미나는 비올 경우 우천 휴장 가능",
        ],
    },
    {
        "title": "아사히카와 여행 (아사히야마 동물원)",
        "description": "삿포로 → 아사히카와 JR 이동. 아사히야마 동물원 + 시내 산책 + 대욕장.",
        "start_date": "2025-11-21",
        "end_date": "2025-11-23",
        "destination": "홋카이도 / 아사히카와",
        "country_code": "JP",
        "total_budget": {"amount_minor": "120000", "currency": "JPY"},
        "days": [
            {
                "region": "삿포로 → 아사히카와",
                "schedules": [
                    ("18:30", "삿포로 → 아사히카와 JR 탑승", "TRANSPORT", "TRAIN", "JR 라일락"),
                    ("19:55", "아사히카와 역 도착", "TRANSPORT", "TRAIN", ""),
                    ("20:00", "호텔 체크인", "ACCOMMODATION", "WALK", ""),
                    ("20:20", "저녁 먹기 + 시내 산책", "MEAL", "WALK", ""),
                    ("22:00", "호텔 대욕장 (15:00–25:00, 5:00–9:00)", "REST", "WALK", ""),
                ],
                "meals": [
                    ("DINNER", "LOCAL", "아사히카와 라멘 또는 징기스칸"),
                ],
                "accom": ("아사히카와 아마넥스 호텔", "20:00", "11:00", "대욕장 / 조식 별도"),
            },
            {
                "region": "아사히카와 (아사히야마 동물원)",
                "schedules": [
                    ("07:30", "기상 + 준비", "REST", "WALK", ""),
                    ("09:40", "역 6번 노리바 → 동물원 버스 (500엔)", "TRANSPORT", "BUS", "9:55→10:34 / 10:10→10:49 / 10:25→11:07"),
                    ("10:22", "아사히야마동물원 도착", "TRANSPORT", "BUS", ""),
                    ("10:30", "동물원 구경 (입장료 1,000엔)", "SIGHTSEEING", "WALK", ""),
                    ("12:00", "동물원 내 점심", "MEAL", "WALK", ""),
                    ("14:10", "동물원 → 아사히카와역 (500엔)", "TRANSPORT", "BUS", "13:40→14:23 / 14:40→15:23 / 15:00→15:45"),
                    ("14:50", "아사히카와역 도착", "TRANSPORT", "BUS", ""),
                    ("15:30", "호텔 대욕장", "REST", "WALK", ""),
                    ("18:00", "시내 구경 + 저녁", "ACTIVITY", "WALK", ""),
                ],
                "meals": [
                    ("LUNCH", "LOCAL", "동물원 식당"),
                    ("DINNER", "LOCAL", "아사히카와 시내 이자카야"),
                ],
                "accom": ("아사히카와 아마넥스 호텔", "—", "11:00", "대욕장"),
            },
            {
                "region": "아사히카와 → 삿포로",
                "schedules": [
                    ("08:00", "기상 + 대욕장 (5:00–9:00)", "REST", "WALK", ""),
                    ("11:00", "체크아웃", "REST", "WALK", ""),
                    ("11:20", "점심 먹기", "MEAL", "WALK", ""),
                    ("12:40", "박물관 / 공원 구경", "SIGHTSEEING", "WALK", ""),
                    ("15:00", "JR 탑승 (아사히카와 → 삿포로)", "TRANSPORT", "TRAIN", "JR 라일락/카무이"),
                    ("16:25", "삿포로 역 도착", "TRANSPORT", "TRAIN", ""),
                    ("17:00", "집 도착", "REST", "WALK", ""),
                ],
                "meals": [
                    ("LUNCH", "LOCAL", "아사히카와 라멘 또는 부타동"),
                ],
                "accom": None,
            },
        ],
        "expenses": [
            ("ACCOMMODATION", 22000, "JPY", "아마넥스 호텔 2박"),
            ("TRANSPORT", 12000, "JPY", "JR 왕복"),
            ("TRANSPORT", 2000, "JPY", "동물원 버스 왕복"),
            ("ACTIVITY", 2000, "JPY", "동물원 입장 2인"),
            ("FOOD", 8000, "JPY", "외식 합계"),
        ],
        "checklist": [
            ("PACKING", "방한 외투 (11월 아사히카와)"),
            ("PACKING", "동물원 사진용 카메라"),
            ("BOOKING", "JR 좌석 지정"),
            ("BOOKING", "아마넥스 호텔 2박"),
            ("TODO", "아사히야마 동물원 펭귄 산책 시간 확인"),
        ],
        "notes": [
            "아마넥스 대욕장 영업: 15:00–25:00 / 05:00–09:00",
            "동물원 노리바 6번, 500엔 버스",
            "11월 후반 적설 가능. 미끄럼 방지 신발 권장",
        ],
    },
    {
        "title": "오비히로 여행 (반에이 경마·부타동)",
        "description": "삿포로 ↔ 오비히로 버스/JR. 오비히로 동물원·미술관·반에이 경마장·부타동·온천.",
        "start_date": "2025-12-27",
        "end_date": "2025-12-30",
        "destination": "홋카이도 / 오비히로",
        "country_code": "JP",
        "total_budget": {"amount_minor": "200000", "currency": "JPY"},
        "days": [
            {
                "region": "삿포로 → 오비히로",
                "schedules": [
                    ("06:30", "일어나기", "REST", "WALK", ""),
                    ("08:10", "삿포로역 4번 노리바로 출발", "TRANSPORT", "WALK", "日本生命札幌ビル / ソラリア西鉄ホテル札幌"),
                    ("09:00", "오비히로행 버스 탑승 + 편의점", "TRANSPORT", "BUS", "포테이토라이너"),
                    ("13:05", "오비히로 버스터미널 도착", "TRANSPORT", "BUS", ""),
                    ("13:30", "호텔에 짐 맡기기", "REST", "WALK", ""),
                    ("14:00", "리치몬드 호텔 오비히로 체크인", "ACCOMMODATION", "WALK", ""),
                    ("15:30", "오비히로 역 근처 산책", "SIGHTSEEING", "WALK", ""),
                    ("18:30", "이자카야 저녁 (오빠 픽)", "MEAL", "WALK", ""),
                ],
                "meals": [
                    ("LUNCH", "LOCAL", "인디아 카레"),
                    ("DINNER", "LOCAL", "이자카야"),
                ],
                "accom": ("리치몬드 호텔 오비히로", "14:00", "11:00", "비즈니스 호텔"),
            },
            {
                "region": "오비히로 (동물원·미술관)",
                "schedules": [
                    ("07:30", "기상 + 준비", "REST", "WALK", ""),
                    ("09:20", "오비히로 버스터미널 2번 노리바 이동", "TRANSPORT", "WALK", ""),
                    ("09:40", "21번 南商業高校線 (200엔)", "TRANSPORT", "BUS", "10:40 / 11:40 / 12:00"),
                    ("09:56", "緑ヶ丘6丁目 帯広美術館入口 하차", "TRANSPORT", "BUS", ""),
                    ("10:00", "오비히로 백년 기념관 (380엔)", "SIGHTSEEING", "WALK", "9시–17시"),
                    ("11:30", "오비히로 동물원 + 점심 (420엔)", "ACTIVITY", "WALK", "11시–14시"),
                    ("15:29", "緑ヶ丘6丁目 → 帯広駅 (200엔)", "TRANSPORT", "BUS", "14:29 / 16:29 / 17:09"),
                    ("15:45", "帯広駅バスターミナル7 도착", "TRANSPORT", "BUS", ""),
                    ("16:00", "온천 즐기기 (히가에리)", "REST", "WALK", "토카치 가덴 1,000엔"),
                    ("18:00", "이자카야 저녁 (오빠 픽)", "MEAL", "WALK", ""),
                ],
                "meals": [
                    ("LUNCH", "LOCAL", "동물원 근처 식당"),
                    ("DINNER", "LOCAL", "이자카야"),
                ],
                "accom": ("리치몬드 호텔 오비히로", "—", "11:00", "비즈니스 호텔"),
            },
            {
                "region": "오비히로 (반에이 경마장)",
                "schedules": [
                    ("09:00", "기상 + 준비", "REST", "WALK", "반에이 경마장 경기일 28~30일"),
                    ("11:00", "호텔 → 버스터미널 12번 노리바", "TRANSPORT", "WALK", ""),
                    ("11:16", "오비히로 경마장행 버스 (240엔)", "TRANSPORT", "BUS", "9:48(31) / 9:54(72) / 10:34(2)"),
                    ("11:30", "반에이 경마장 도착 (10시 오픈)", "SIGHTSEEING", "WALK", "토카치무라 점심"),
                    ("14:00", "1R 경마 관람", "ACTIVITY", "WALK", ""),
                    ("15:02", "경마장 → 버스터미널 (240엔)", "TRANSPORT", "BUS", "13:02(1) / 13:37(72) / 13:40(17)"),
                    ("15:17", "오비히로역 도착", "TRANSPORT", "BUS", ""),
                    ("16:00", "온천 즐기기", "REST", "WALK", "프레미아 캐빈 1,200엔"),
                    ("18:30", "이자카야 저녁 (오빠 픽)", "MEAL", "WALK", ""),
                ],
                "meals": [
                    ("LUNCH", "LOCAL", "토카치무라 (경마장 내)"),
                    ("DINNER", "LOCAL", "이자카야"),
                ],
                "accom": ("리치몬드 호텔 오비히로", "—", "11:00", "비즈니스 호텔"),
            },
            {
                "region": "오비히로 → 삿포로",
                "schedules": [
                    ("09:00", "기상 + 준비", "REST", "WALK", ""),
                    ("11:00", "체크아웃", "REST", "WALK", ""),
                    ("11:30", "부타동 먹기 (豚丼のぶたはげ 10시 오픈)", "MEAL", "WALK", ""),
                    ("13:12", "오비히로 → 삿포로 JR 탑승", "TRANSPORT", "TRAIN", "JR 오조라"),
                    ("16:10", "삿포로 역 도착", "TRANSPORT", "TRAIN", ""),
                    ("16:30", "집 도착", "REST", "WALK", ""),
                ],
                "meals": [
                    ("LUNCH", "LOCAL", "豚丼のぶたはげ"),
                ],
                "accom": None,
            },
        ],
        "expenses": [
            ("ACCOMMODATION", 30000, "JPY", "리치몬드 호텔 3박"),
            ("TRANSPORT", 7600, "JPY", "포테이토라이너 + JR"),
            ("TRANSPORT", 2000, "JPY", "시내 버스 합계"),
            ("ACTIVITY", 1600, "JPY", "동물원 + 백년기념관"),
            ("ACTIVITY", 2200, "JPY", "히가에리 온천 2회"),
            ("FOOD", 12000, "JPY", "외식 합계"),
        ],
        "checklist": [
            ("PACKING", "방한 장비 (12월말 오비히로)"),
            ("PACKING", "현금 (버스 잔돈)"),
            ("BOOKING", "포테이토라이너 좌석"),
            ("BOOKING", "리치몬드 호텔 3박"),
            ("BOOKING", "JR 오조라 좌석 지정"),
            ("TODO", "반에이 경마 1R 시간 14:00 확인"),
            ("TODO", "豚丼のぶたはげ 10시 오픈"),
        ],
        "notes": [
            "히가에리 온천 후보: 토카치 가덴 1,000엔 / 프레미아 캐빈 1,200엔 / 후쿠이 호텔 1,500엔",
            "후지마루 파크는 올해 영업 종료",
            "반에이 경마장 경기일: 12/28~30",
        ],
    },
]


def seed_one(token: str, t: dict[str, Any]) -> str:
    print(f"→ trip 생성: {t['title']}")
    res = api(token, "POST", "/api/v1/trips", {
        "title": t["title"],
        "description": t["description"],
        "startDate": t["start_date"],
        "endDate": t["end_date"],
        "destination": t["destination"],
        "countryCode": t["country_code"],
        "totalBudget": t["total_budget"],
    })
    trip = res["trip"]
    trip_id = trip["id"]
    print(f"   tripId={trip_id}")

    # 자동 생성된 days 조회.
    days_res = api(token, "GET", f"/api/v1/trips/{trip_id}/days")
    days = days_res.get("days", [])
    if len(days) != len(t["days"]):
        print(f"   ! day 수 불일치: api={len(days)} pdf={len(t['days'])}")

    for i, (api_day, src_day) in enumerate(zip(days, t["days"], strict=False)):
        day_id = api_day["id"]
        # region 업데이트.
        api(token, "PATCH", f"/api/v1/days/{day_id}", {"region": src_day["region"]})

        # schedules.
        for st, title, cat, transport, detail in src_day["schedules"]:
            api(token, "POST", f"/api/v1/days/{day_id}/schedules", {
                "startTime": st,
                "title": title,
                "region": src_day["region"],
                "category": f"SCHEDULE_CATEGORY_{cat}",
                "transport": f"TRANSPORT_TYPE_{transport}",
                "transportDetail": detail,
            })

        # meals.
        for mt, src, name in src_day["meals"]:
            api(token, "PUT", f"/api/v1/days/{day_id}/meals", {
                "mealType": f"MEAL_TYPE_{mt}",
                "source": f"MEAL_SOURCE_{src}",
                "restaurantName": name,
            })

        # accommodation.
        if src_day["accom"]:
            name, ci, co, amen = src_day["accom"]
            api(token, "PUT", f"/api/v1/days/{day_id}/accommodation", {
                "name": name,
                "checkInTime": ci,
                "checkOutTime": co,
                "amenities": amen,
            })
        print(f"   ✓ day {i + 1}/{len(days)}")

    # expenses.
    for cat, amount, cur, desc in t.get("expenses", []):
        try:
            api(token, "POST", f"/api/v1/trips/{trip_id}/expenses", {
                "category": f"EXPENSE_CATEGORY_{cat}",
                "amount": {"amountMinor": str(amount * 100), "currency": cur},
                "description": desc,
                "occurredAt": t["start_date"] + "T12:00:00Z",
            })
        except RuntimeError as e:
            print(f"   ! expense 실패: {e}")

    # checklist.
    for cat, item in t.get("checklist", []):
        try:
            api(token, "POST", f"/api/v1/trips/{trip_id}/checklist", {
                "category": f"CHECKLIST_CATEGORY_{cat}",
                "title": item,
            })
        except RuntimeError as e:
            print(f"   ! checklist 실패: {e}")

    # notes.
    for content in t.get("notes", []):
        try:
            api(token, "POST", f"/api/v1/trips/{trip_id}/notes", {
                "content": content,
            })
        except RuntimeError as e:
            print(f"   ! note 실패: {e}")

    return trip_id


def main() -> None:
    token = get_token()
    print(f"✓ Keycloak token 획득 (len={len(token)})")

    # 기존 trip 정리.
    if os.environ.get("CLEAN", "1") == "1":
        existing = api(token, "GET", "/api/v1/trips")
        for tr in existing.get("trips", []):
            print(f"× 기존 trip 삭제: {tr['title']} ({tr['id']})")
            try:
                api(token, "DELETE", f"/api/v1/trips/{tr['id']}")
            except RuntimeError as e:
                print(f"   ! 삭제 실패: {e}")

    ids = []
    for t in TRIPS:
        try:
            ids.append(seed_one(token, t))
            time.sleep(0.3)
        except Exception as e:
            print(f"!! {t['title']} 시드 실패: {e}")

    print(f"\n✓ 완료: {len(ids)}개 trip 시드 (testuser)")
    for i, tid in enumerate(ids):
        print(f"  - {TRIPS[i]['title']}: {tid}")


if __name__ == "__main__":
    main()
