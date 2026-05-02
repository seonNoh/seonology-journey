# Seonology Journey Android

Compose + AppAuth-Android (Keycloak PKCE) + Retrofit/Moshi + Coroutines 기반의 모바일 클라이언트 스켈레톤.

## 빌드 환경

- Android Studio Ladybug (2024.2.1) 이상
- JDK 17
- Android SDK 35 (compile/target), minSdk 26
- Kotlin 2.0.21, AGP 8.7.2

## 실행

```sh
cd apps/android
./gradlew :app:assembleDebug
```

루트에 gradle wrapper 가 없는 경우 `gradle wrapper --gradle-version=8.10.2` 로 생성.

## BuildConfig

`app/build.gradle.kts` 에 정의:

| 키 | 기본값 |
| --- | --- |
| `API_BASE` | `https://journey-api.seonology.com` |
| `KEYCLOAK_ISSUER` | `https://keycloak.seonology.com/realms/seonology` |
| `KEYCLOAK_CLIENT_ID` | `seonology-journey-android` |

AppAuth 콜백 URI: `com.seonology.journey:/oauth2redirect`
Keycloak 클라이언트의 Valid Redirect URIs 에 위 값 추가.

## 주요 모듈

- `auth/` — `AuthStore` (SharedPreferences), `KeycloakAuth` (PKCE), `AuthInterceptor`
- `data/` — Retrofit `JourneyApi`, Moshi 모델
- `ui/JourneyApp.kt` — Compose Scaffold + Trip 목록 화면

## 향후

- Trip 상세 / Day 상세 / Schedule 편집 / Media 업로드 화면
- WebSocket 실시간 동기화 (`/ws/trips/{tripId}`)
- EncryptedSharedPreferences 또는 Tink 로 토큰 보호
- Hilt DI / Room offline cache
