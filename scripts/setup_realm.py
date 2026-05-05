#!/usr/bin/env python3
"""Create seonology-journey realm with clients and users."""
import json, sys, urllib.request, urllib.parse

BASE = "https://auth.seonology.com"

def post(url, data=None, token=None, method="POST"):
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    body = json.dumps(data).encode() if data else None
    req = urllib.request.Request(url, data=body, headers=headers, method=method)
    try:
        resp = urllib.request.urlopen(req)
        ct = resp.headers.get("Content-Type", "")
        if "json" in ct:
            return resp.status, json.loads(resp.read())
        return resp.status, resp.read().decode()
    except urllib.error.HTTPError as e:
        body = e.read().decode()
        return e.code, body

def form_post(url, data):
    body = urllib.parse.urlencode(data).encode()
    req = urllib.request.Request(url, data=body,
        headers={"Content-Type": "application/x-www-form-urlencoded"})
    resp = urllib.request.urlopen(req)
    return json.loads(resp.read())

# 1. Admin token
print("1. Admin token 획득...")
tok = form_post(f"{BASE}/realms/master/protocol/openid-connect/token", {
    "grant_type": "password", "client_id": "admin-cli",
    "username": "admin", "password": "keycloak-admin-pass"
})
AT = tok["access_token"]
print(f"   token len={len(AT)}")

# 2. Create realm
print("2. Realm seonology-journey 생성...")
code, resp = post(f"{BASE}/admin/realms", {
    "realm": "seonology-journey",
    "enabled": True,
    "registrationAllowed": False,
    "loginWithEmailAllowed": True,
    "duplicateEmailsAllowed": False,
    "sslRequired": "external",
    "accessTokenLifespan": 300,
    "ssoSessionIdleTimeout": 1800,
    "ssoSessionMaxLifespan": 36000,
}, AT)
print(f"   → {code} {resp[:200] if isinstance(resp, str) else resp}")

# 3. journey-web client
print("3. journey-web client 생성...")
code, resp = post(f"{BASE}/admin/realms/seonology-journey/clients", {
    "clientId": "journey-web",
    "enabled": True,
    "publicClient": True,
    "directAccessGrantsEnabled": True,
    "standardFlowEnabled": True,
    "rootUrl": "https://journey.seonology.com",
    "redirectUris": [
        "https://journey.seonology.com/*",
        "http://localhost:5173/*"
    ],
    "webOrigins": [
        "https://journey.seonology.com",
        "http://localhost:5173"
    ],
    "protocol": "openid-connect",
}, AT)
print(f"   → {code}")

# 4. journey-android client
print("4. journey-android client 생성...")
code, resp = post(f"{BASE}/admin/realms/seonology-journey/clients", {
    "clientId": "journey-android",
    "enabled": True,
    "publicClient": True,
    "directAccessGrantsEnabled": False,
    "standardFlowEnabled": True,
    "redirectUris": ["com.seonology.journey://oauth2redirect"],
    "protocol": "openid-connect",
}, AT)
print(f"   → {code}")

# 5. Create users
USERS = [
    {"username": "seon", "password": "!Tjs78xor0512", "email": "seon@seonology.com", "firstName": "Seon", "lastName": "Noh"},
    {"username": "kkogu", "password": "Akcmwns0830", "email": "kkogu@seonology.com", "firstName": "Kkogu", "lastName": "User"},
]

for u in USERS:
    print(f"5. user '{u['username']}' 생성...")
    code, resp = post(f"{BASE}/admin/realms/seonology-journey/users", {
        "username": u["username"],
        "enabled": True,
        "emailVerified": True,
        "email": u["email"],
        "firstName": u["firstName"],
        "lastName": u["lastName"],
        "credentials": [{
            "type": "password",
            "value": u["password"],
            "temporary": False,
        }],
    }, AT)
    print(f"   → {code}")

# 6. Verify: list users
print("6. Realm 유저 목록 확인...")
code, users = post(f"{BASE}/admin/realms/seonology-journey/users?max=50", token=AT, method="GET", data=None)
# GET doesn't use data
req = urllib.request.Request(f"{BASE}/admin/realms/seonology-journey/users?max=50",
    headers={"Authorization": f"Bearer {AT}"})
resp = urllib.request.urlopen(req)
users = json.loads(resp.read())
for u in users:
    print(f"   - {u['username']} (id={u['id']})")

print("\n✓ realm seonology-journey 설정 완료")
