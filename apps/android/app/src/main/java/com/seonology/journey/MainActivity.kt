package com.seonology.journey

// App entry point. 릴리즈 CI 가 이 모듈을 gradle assembleRelease 로 빌드하면
// Firebase App Distribution 이 internal 그룹 테스터에게 APK 를 배포한다.
//
// 알림 없이 설치가 바로 되게 하려면 테스터는 Firebase App Tester 앱에서
// 사전에 로그인 상태를 유지해 두면 된다.
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.Surface
import androidx.compose.material3.MaterialTheme
import androidx.compose.ui.Modifier
import com.seonology.journey.ui.JourneyApp
import com.seonology.journey.ui.theme.JourneyTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        // Edge-to-edge so the sakura gradient background extends under status bar.
        enableEdgeToEdge()
        setContent {
            JourneyTheme {
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background,
                ) {
                    JourneyApp()
                }
            }
        }
    }
}
