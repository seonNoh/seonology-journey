package com.seonology.journey.auth

import android.content.Context
import androidx.core.content.edit

/**
 * 단순 SharedPreferences 기반 토큰 저장소.
 * 실제 운영에서는 EncryptedSharedPreferences 또는 Tink 사용을 권장.
 */
class AuthStore(context: Context) {
    private val prefs = context.getSharedPreferences("journey_auth", Context.MODE_PRIVATE)

    var accessToken: String?
        get() = prefs.getString(KEY_ACCESS, null)
        set(value) = prefs.edit { putString(KEY_ACCESS, value) }

    var refreshToken: String?
        get() = prefs.getString(KEY_REFRESH, null)
        set(value) = prefs.edit { putString(KEY_REFRESH, value) }

    val isAuthenticated: Boolean get() = !accessToken.isNullOrEmpty()

    fun clear() {
        prefs.edit { clear() }
    }

    companion object {
        private const val KEY_ACCESS = "access_token"
        private const val KEY_REFRESH = "refresh_token"
    }
}
