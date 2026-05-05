"""
Seed DynamoDB with travel data parsed from source PDFs.
Usage: AWS_PROFILE=seonology python3 scripts/seed_from_pdf.py
"""
import boto3
import uuid
from datetime import datetime

REGION = "ap-northeast-1"
TABLE_PREFIX = "seonology-journey-"
OWNER_ID = "32e80a79-21fa-4039-9c93-55f09fb19bda"  # seon (Keycloak)
NOW = datetime.utcnow().strftime("%Y-%m-%dT%H:%M:%SZ")


def get_client():
    return boto3.client("dynamodb", region_name=REGION)


def put_trip(c, trip_id, title, description, start_date, end_date, status="completed"):
    c.put_item(
        TableName=f"{TABLE_PREFIX}trips",
        Item={
            "PK": {"S": f"USER#{OWNER_ID}"},
            "SK": {"S": f"TRIP#{trip_id}"},
            "GSI1PK": {"S": f"TRIP#{trip_id}"},
            "GSI1SK": {"S": "METADATA"},
            "id": {"S": trip_id},
            "ownerId": {"S": OWNER_ID},
            "userId": {"S": OWNER_ID},
            "title": {"S": title},
            "description": {"S": description},
            "startDate": {"S": start_date},
            "endDate": {"S": end_date},
            "status": {"S": status},
            "visibility": {"S": "private"},
            "totalBudget": {"N": "0"},
            "createdAt": {"S": NOW},
            "updatedAt": {"S": NOW},
        },
    )
    print(f"  [trip] {title} ({start_date} ~ {end_date})")


def put_day(c, trip_id, day_id, day_number, date, region=""):
    c.put_item(
        TableName=f"{TABLE_PREFIX}days",
        Item={
            "PK": {"S": f"TRIP#{trip_id}"},
            "SK": {"S": f"DAY#{day_number:03d}"},
            "id": {"S": day_id},
            "dayNumber": {"N": str(day_number)},
            "date": {"S": date},
            "region": {"S": region},
            "createdAt": {"S": NOW},
            "updatedAt": {"S": NOW},
        },
    )
    print(f"    [day] {day_number} - {date} ({region})")
    return day_id


def put_schedule(c, day_id, sort_order, title, start_time="", category="activity", location="", transport="", notes=""):
    schedule_id = str(uuid.uuid4())
    item = {
        "PK": {"S": f"DAY#{day_id}"},
        "SK": {"S": f"SCHEDULE#{sort_order:03d}"},
        "id": {"S": schedule_id},
        "title": {"S": title},
        "sortOrder": {"N": str(sort_order)},
        "category": {"S": category},
        "createdAt": {"S": NOW},
        "updatedAt": {"S": NOW},
    }
    if start_time:
        item["startTime"] = {"S": start_time}
    if location:
        item["location"] = {"S": location}
    if transport:
        item["transport"] = {"S": transport}
    if notes:
        item["notes"] = {"S": notes}
    c.put_item(TableName=f"{TABLE_PREFIX}schedules", Item=item)
    print(f"      [schedule] {sort_order}: {start_time} {title}")


def put_meal(c, day_id, meal_type, restaurant="", source="local", rating=0, notes=""):
    meal_id = str(uuid.uuid4())
    item = {
        "PK": {"S": f"DAY#{day_id}"},
        "SK": {"S": f"MEAL#{meal_type}"},
        "id": {"S": meal_id},
        "mealType": {"S": meal_type},
        "restaurant": {"S": restaurant},
        "source": {"S": source},
        "rating": {"N": str(rating)},
        "createdAt": {"S": NOW},
        "updatedAt": {"S": NOW},
    }
    if notes:
        item["notes"] = {"S": notes}
    c.put_item(TableName=f"{TABLE_PREFIX}meals", Item=item)
    print(f"      [meal] {meal_type}: {restaurant or source}")


def put_accommodation(c, day_id, name, check_in="", check_out="", notes=""):
    acc_id = str(uuid.uuid4())
    item = {
        "PK": {"S": f"DAY#{day_id}"},
        "SK": {"S": "ACCOMMODATION"},
        "AccommodationId": {"S": acc_id},
        "id": {"S": acc_id},
        "name": {"S": name},
        "createdAt": {"S": NOW},
        "updatedAt": {"S": NOW},
    }
    if check_in:
        item["checkIn"] = {"S": check_in}
    if check_out:
        item["checkOut"] = {"S": check_out}
    if notes:
        item["notes"] = {"S": notes}
    c.put_item(TableName=f"{TABLE_PREFIX}accommodations", Item=item)
    print(f"      [accommodation] {name}")


def seed_akan(c):
    """아칸 여행 (2025-10-09 ~ 2025-10-12) - 기존 trip 재활용, schedule 추가"""
    print("\n=== 아칸 여행 ===")
    trip_id = "8a1b2c3d-4e5f-6789-abcd-ef0123456789"

    # 기존 Day IDs
    days = {
        1: "cd6b115e-8a15-4f3a-ae66-b4f37ff62e2c",  # 10/9
        2: "b3dfe793-b032-4a91-addc-52a1c18f40d2",  # 10/10
        3: "ac7eab58-a55f-4ff7-b3b9-df5cae063bcb",  # 10/11
        4: "e720613b-08ad-428c-9909-457038f13019",  # 10/12
    }

    # Day 1: 삿포로 -> 치토세 -> 쿠시로 -> 아칸
    d1 = days[1]
    put_schedule(c, d1, 1, "일어나기", "08:00", "activity", "삿포로")
    put_schedule(c, d1, 2, "삿포로역 출발", "09:30", "transport", "삿포로역", "train")
    put_schedule(c, d1, 3, "에어포트 탑승", "10:00", "transport", "삿포로역", "train", "짐 붙이기")
    put_schedule(c, d1, 4, "신치토세공항 도착", "11:00", "transport", "신치토세공항")
    put_schedule(c, d1, 5, "점심 & 리락쿠마/쵸파 찾아보기", "11:30", "activity", "신치토세공항")
    put_schedule(c, d1, 6, "ANA4873 신치토세공항 출발", "13:20", "transport", "신치토세공항", "flight", "ANA4873")
    put_schedule(c, d1, 7, "쿠시로 공항 도착", "14:05", "transport", "쿠시로공항", "flight")
    put_schedule(c, d1, 8, "버스 탑승 -> 아칸 이동", "15:37", "transport", "쿠시로->아칸", "bus")
    put_schedule(c, d1, 9, "호텔 체크인", "17:00", "accommodation", "아칸")
    put_schedule(c, d1, 10, "저녁 식사 & 마리모 마츠리", "18:00", "activity", "아칸")
    put_schedule(c, d1, 11, "온천", "20:00", "activity", "호텔")
    put_meal(c, d1, "lunch", "신치토세공항 현지식", "local")
    put_meal(c, d1, "dinner", "호텔식", "hotel")
    put_accommodation(c, d1, "아칸츠루가별장 히나노자 (あかん鶴雅別荘 鄙の座)", "17:00", "11:00")

    # Day 2: 아칸 관광
    d2 = days[2]
    put_schedule(c, d2, 1, "호텔 조식", "07:30", "meal", "호텔")
    put_schedule(c, d2, 2, "온천 즐기기", "08:30", "activity", "호텔")
    put_schedule(c, d2, 3, "마리모 마츠리 구경", "10:00", "activity", "아칸호", "10:00/10:30/11:30")
    put_schedule(c, d2, 4, "阿寒湖畔エコミュージアムセンター & ボッケ遊歩道", "11:00", "activity", "아칸호반 에코뮤지엄센터", notes="약 1시간 소요")
    put_schedule(c, d2, 5, "유람선 85분", "11:00", "activity", "아칸호", notes="1인 2,400엔 (마리모 관람 포함). 9:00~16:00 매시 출발")
    put_schedule(c, d2, 6, "점심", "13:00", "meal", "아칸호 주변")
    put_schedule(c, d2, 7, "阿寒湖アイヌコタン", "14:00", "activity", "아칸호 아이누코탄")
    put_schedule(c, d2, 8, "호텔 복귀 & 온천 휴식", "16:30", "activity", "호텔")
    put_schedule(c, d2, 9, "저녁 식사", "17:30", "meal", "호텔")
    put_schedule(c, d2, 10, "카무이 루미나 관람", "18:00", "activity", "아칸호", notes="약 1시간. 17:30~21:00 15분 단위. 호텔에서 티켓 구매")
    put_meal(c, d2, "breakfast", "호텔식", "hotel")
    put_meal(c, d2, "lunch", "현지식", "local")
    put_meal(c, d2, "dinner", "호텔식", "hotel")
    put_accommodation(c, d2, "아칸츠루가별장 히나노자 (あかん鶴雅別荘 鄙の座)", "17:00", "11:00")

    # Day 3: 아칸 -> 쿠시로
    d3 = days[3]
    put_schedule(c, d3, 1, "호텔 조식", "07:30", "meal", "호텔")
    put_schedule(c, d3, 2, "온천 즐기기", "08:30", "activity", "호텔")
    put_schedule(c, d3, 3, "체크아웃 (짐 붙이기)", "11:00", "activity", "호텔", notes="체크아웃 11시까지")
    put_schedule(c, d3, 4, "온천 마을 산책", "11:30", "activity", "아칸 온천마을")
    put_schedule(c, d3, 5, "버스 탑승 -> 쿠시로 이동", "12:40", "transport", "아칸->쿠시로", "bus", "약 2시간, 2,750엔")
    put_schedule(c, d3, 6, "쿠시로역 도착", "14:40", "transport", "쿠시로역")
    put_schedule(c, d3, 7, "호텔 체크인", "15:00", "accommodation", "쿠시로 스파 호텔")
    put_schedule(c, d3, 8, "쿠시로역 근처 산책", "15:30", "activity", "쿠시로역 주변")
    put_schedule(c, d3, 9, "저녁 식사", "17:30", "meal", "쿠시로")
    put_schedule(c, d3, 10, "호텔 휴식", "19:00", "activity", "호텔")
    put_meal(c, d3, "breakfast", "호텔식", "hotel")
    put_meal(c, d3, "lunch", "현지식", "local")
    put_meal(c, d3, "dinner", "현지식 (치킨)", "local")
    put_accommodation(c, d3, "쿠시로 스파 호텔", "15:00", "10:00")

    # Day 4: 쿠시로 -> 삿포로
    d4 = days[4]
    # 기존 schedule 이 있을 수 있으므로 덮어쓰기
    put_schedule(c, d4, 1, "호텔 조식 없음", "08:00", "activity", "호텔")
    put_schedule(c, d4, 2, "체크아웃", "10:00", "activity", "호텔", notes="체크아웃 10시까지")
    put_schedule(c, d4, 3, "버스 탑승 -> 쿠시로공항", "12:15", "transport", "쿠시로->쿠시로공항", "bus", "連絡バス 950엔")
    put_schedule(c, d4, 4, "쿠시로공항 도착", "13:00", "transport", "쿠시로공항")
    put_schedule(c, d4, 5, "ANA4874 쿠시로공항 출발", "14:35", "transport", "쿠시로공항", "flight", "ANA4874")
    put_schedule(c, d4, 6, "신치토세공항 도착", "15:20", "transport", "신치토세공항", "flight")
    put_schedule(c, d4, 7, "에어포트 탑승", "16:00", "transport", "신치토세공항", "train")
    put_schedule(c, d4, 8, "삿포로역 도착", "17:00", "transport", "삿포로역", "train")
    put_meal(c, d4, "lunch", "현지식", "local", notes="공항에서")
    put_accommodation(c, d4, "자택", "", "")


def seed_asahikawa(c):
    """아사히카와 여행 (2025-11-21 ~ 2025-11-23)"""
    print("\n=== 아사히카와 여행 ===")
    trip_id = str(uuid.uuid4())
    put_trip(c, trip_id, "아사히카와 여행", "2025년 11월 아사히야마 동물원 & 아사히카와 시내 관광", "2025-11-21", "2025-11-23")

    # Day 1: 삿포로 -> 아사히카와
    d1 = str(uuid.uuid4())
    put_day(c, trip_id, d1, 1, "2025-11-21", "아사히카와")
    put_schedule(c, d1, 1, "아사히카와행 JR 탑승", "18:30", "transport", "삿포로역", "train")
    put_schedule(c, d1, 2, "아사히카와역 도착", "19:55", "transport", "아사히카와역", "train")
    put_schedule(c, d1, 3, "호텔 체크인", "20:00", "accommodation", "아사히카와 아마넥스 호텔")
    put_schedule(c, d1, 4, "저녁 식사 & 시내 구경", "20:20", "meal", "아사히카와 시내")
    put_schedule(c, d1, 5, "호텔 복귀 & 대욕장", "22:00", "activity", "호텔", notes="대욕장 15:00~25:00, 5:00~9:00")
    put_meal(c, d1, "dinner", "현지식", "local")
    put_accommodation(c, d1, "아사히카와 아마넥스 호텔", "20:00", "11:00")

    # Day 2: 아사히야마 동물원
    d2 = str(uuid.uuid4())
    put_day(c, trip_id, d2, 2, "2025-11-22", "아사히카와")
    put_schedule(c, d2, 1, "기상 & 준비", "07:30", "activity", "호텔")
    put_schedule(c, d2, 2, "아사히야마 동물원행 버스 탑승", "09:40", "transport", "아사히카와역 6번 노리바", "bus", "500엔")
    put_schedule(c, d2, 3, "아사히야마 동물원 도착", "10:22", "transport", "아사히야마동물원", "bus")
    put_schedule(c, d2, 4, "동물원 구경", "10:30", "activity", "아사히야마동물원", notes="입장료 1,000엔")
    put_schedule(c, d2, 5, "점심", "12:00", "meal", "아사히야마동물원")
    put_schedule(c, d2, 6, "아사히카와역행 버스 탑승", "14:10", "transport", "아사히야마동물원", "bus", "500엔")
    put_schedule(c, d2, 7, "아사히카와역 도착", "14:50", "transport", "아사히카와역", "bus")
    put_schedule(c, d2, 8, "호텔 대욕장", "15:30", "activity", "호텔")
    put_schedule(c, d2, 9, "아사히카와 시내 구경 & 저녁", "17:00", "activity", "아사히카와 시내")
    put_meal(c, d2, "lunch", "현지식", "local")
    put_meal(c, d2, "dinner", "현지식", "local")
    put_accommodation(c, d2, "아사히카와 아마넥스 호텔", "20:00", "11:00")

    # Day 3: 아사히카와 -> 삿포로
    d3 = str(uuid.uuid4())
    put_day(c, trip_id, d3, 3, "2025-11-23", "아사히카와")
    put_schedule(c, d3, 1, "기상 & 대욕장", "07:00", "activity", "호텔", notes="대욕장 5:00~9:00")
    put_schedule(c, d3, 2, "체크아웃", "11:00", "activity", "호텔")
    put_schedule(c, d3, 3, "점심 식사", "11:20", "meal", "아사히카와 시내")
    put_schedule(c, d3, 4, "아사히카와 뒤쪽 구경 (박물원, 공원)", "12:40", "activity", "아사히카와")
    put_schedule(c, d3, 5, "JR 탑승 -> 삿포로", "15:00", "transport", "아사히카와역", "train")
    put_schedule(c, d3, 6, "삿포로역 도착", "16:25", "transport", "삿포로역", "train")
    put_schedule(c, d3, 7, "집 도착", "17:00", "activity", "삿포로")
    put_meal(c, d3, "lunch", "현지식", "local")
    put_accommodation(c, d3, "자택", "", "")


def seed_obihiro(c):
    """오비히로 여행 (2025-12-27 ~ 2025-12-30)"""
    print("\n=== 오비히로 여행 ===")
    trip_id = str(uuid.uuid4())
    put_trip(c, trip_id, "오비히로 여행", "2025년 12월 오비히로 동물원, 반에이경마장, 온천 여행", "2025-12-27", "2025-12-30")

    # Day 1: 삿포로 -> 오비히로
    d1 = str(uuid.uuid4())
    put_day(c, trip_id, d1, 1, "2025-12-27", "오비히로")
    put_schedule(c, d1, 1, "일어나기", "06:30", "activity", "삿포로")
    put_schedule(c, d1, 2, "버스 터미널 출발", "08:10", "transport", "삿포로역앞 4번", notes="日本生命札幌ビル, ソラリア西鉄ホテル札幌")
    put_schedule(c, d1, 3, "오비히로행 버스 탑승", "09:00", "transport", "삿포로", "bus", "편의점에서 음료/간식 구매")
    put_schedule(c, d1, 4, "오비히로 버스터미널 도착", "13:05", "transport", "오비히로 버스터미널", "bus")
    put_schedule(c, d1, 5, "호텔에 짐 맡기기", "13:20", "activity", "리치몬드 호텔 오비히로")
    put_schedule(c, d1, 6, "점심 (인디아 카레)", "13:30", "meal", "오비히로 시내")
    put_schedule(c, d1, 7, "호텔 체크인", "14:00", "accommodation", "리치몬드 호텔 오비히로", notes="체크인 14시부터")
    put_schedule(c, d1, 8, "온천 또는 오비히로역 근처 산책", "15:00", "activity", "오비히로역 주변")
    put_schedule(c, d1, 9, "저녁 (이자카야)", "18:00", "meal", "오비히로 시내")
    put_schedule(c, d1, 10, "호텔 복귀 & 휴식", "20:00", "activity", "호텔")
    put_meal(c, d1, "lunch", "인디아 카레", "local")
    put_meal(c, d1, "dinner", "이자카야", "local")
    put_accommodation(c, d1, "리치몬드 호텔 오비히로", "14:00", "11:00")

    # Day 2: 오비히로 백년기념관 + 동물원
    d2 = str(uuid.uuid4())
    put_day(c, trip_id, d2, 2, "2025-12-28", "오비히로")
    put_schedule(c, d2, 1, "기상 & 준비", "07:30", "activity", "호텔")
    put_schedule(c, d2, 2, "오비히로 버스터미널 2번 이동", "09:20", "transport", "오비히로 버스터미널")
    put_schedule(c, d2, 3, "버스 탑승 (21 南商業高校線)", "09:40", "transport", "오비히로 버스터미널 2번", "bus", "200엔")
    put_schedule(c, d2, 4, "緑ヶ丘6丁目 하차 (帯広美術館入口)", "09:56", "transport", "緑ヶ丘6丁目", "bus", "하차 후 도보 10분")
    put_schedule(c, d2, 5, "오비히로 백년기념관", "10:10", "activity", "오비히로 백년기념관", notes="입장료 380엔, 9시~17시")
    put_schedule(c, d2, 6, "동물원 구경 & 점심", "11:00", "activity", "오비히로동물원", notes="입장료 420엔, 11시~14시")
    put_schedule(c, d2, 7, "버스 탑승 -> 역 복귀", "15:29", "transport", "帯広美術館入口", "bus", "200엔 (21 緑ヶ丘6丁目)")
    put_schedule(c, d2, 8, "오비히로역 도착", "15:45", "transport", "帯広駅バスターミナル7", "bus")
    put_schedule(c, d2, 9, "온천 즐기기", "16:00", "activity", "오비히로 온천")
    put_schedule(c, d2, 10, "호텔에서 휴식", "17:00", "activity", "호텔")
    put_schedule(c, d2, 11, "저녁 (이자카야)", "18:00", "meal", "오비히로 시내")
    put_meal(c, d2, "lunch", "현지식 (동물원 근처)", "local")
    put_meal(c, d2, "dinner", "이자카야", "local")
    put_accommodation(c, d2, "리치몬드 호텔 오비히로", "14:00", "11:00")

    # Day 3: 반에이경마장
    d3 = str(uuid.uuid4())
    put_day(c, trip_id, d3, 3, "2025-12-29", "오비히로")
    put_schedule(c, d3, 1, "기상 & 준비", "09:00", "activity", "호텔")
    put_schedule(c, d3, 2, "호텔 나가기 -> 버스터미널 12번", "11:00", "transport", "오비히로 버스터미널")
    put_schedule(c, d3, 3, "버스 탑승 -> 경마장", "11:16", "transport", "오비히로 버스터미널", "bus", "240엔")
    put_schedule(c, d3, 4, "경마장 도착", "11:30", "transport", "오비히로 경마장", "bus")
    put_schedule(c, d3, 5, "토카치무라에서 점심", "11:45", "meal", "토카치무라 (경마장내)", notes="대부분 11시 오픈")
    put_schedule(c, d3, 6, "반에이경마장 구경", "13:00", "activity", "반에이경마장", notes="10시 오픈, 1R 14시, 경기일 28~30일")
    put_schedule(c, d3, 7, "버스 탑승 -> 역 복귀", "15:02", "transport", "오비히로경마장앞", "bus", "240엔")
    put_schedule(c, d3, 8, "오비히로역 도착", "15:17", "transport", "오비히로역", "bus")
    put_schedule(c, d3, 9, "온천 즐기기", "16:00", "activity", "오비히로 온천")
    put_schedule(c, d3, 10, "호텔 휴식", "17:00", "activity", "호텔")
    put_schedule(c, d3, 11, "저녁 (이자카야)", "18:30", "meal", "오비히로 시내")
    put_meal(c, d3, "lunch", "토카치무라 (경마장)", "local")
    put_meal(c, d3, "dinner", "이자카야", "local")
    put_accommodation(c, d3, "리치몬드 호텔 오비히로", "14:00", "11:00")

    # Day 4: 오비히로 -> 삿포로
    d4 = str(uuid.uuid4())
    put_day(c, trip_id, d4, 4, "2025-12-30", "오비히로")
    put_schedule(c, d4, 1, "기상 & 준비", "09:00", "activity", "호텔")
    put_schedule(c, d4, 2, "체크아웃", "11:00", "activity", "호텔")
    put_schedule(c, d4, 3, "부타동 점심", "11:30", "meal", "오비히로역", notes="豚丼のぶたはげ 10시 오픈")
    put_schedule(c, d4, 4, "JR 탑승 -> 삿포로", "13:12", "transport", "오비히로역", "train")
    put_schedule(c, d4, 5, "삿포로역 도착", "16:10", "transport", "삿포로역", "train")
    put_schedule(c, d4, 6, "집 도착", "16:30", "activity", "삿포로")
    put_meal(c, d4, "lunch", "부타동 (豚丼のぶたはげ)", "local")
    put_accommodation(c, d4, "자택", "", "")


def main():
    c = get_client()
    print("PDF 여행 데이터 시드 시작...")

    seed_akan(c)
    seed_asahikawa(c)
    seed_obihiro(c)

    print("\n=== 시드 완료 ===")


if __name__ == "__main__":
    main()
