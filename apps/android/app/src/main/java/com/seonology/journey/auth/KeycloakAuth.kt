package com.seonology.journey.auth

import android.app.Activity
import android.content.Context
import android.net.Uri
import com.seonology.journey.BuildConfig
import net.openid.appauth.AuthorizationRequest
import net.openid.appauth.AuthorizationService
import net.openid.appauth.AuthorizationServiceConfiguration
import net.openid.appauth.ResponseTypeValues
import net.openid.appauth.TokenResponse

/**
 * Keycloak PKCE 흐름 헬퍼.
 *
 * 사용 예:
 *   val cfg = KeycloakAuth.config()
 *   val service = AuthorizationService(activity)
 *   val intent = service.getAuthorizationRequestIntent(KeycloakAuth.buildAuthRequest(cfg))
 *   activity.startActivityForResult(intent, RC_AUTH)
 */
object KeycloakAuth {
    private const val REDIRECT_URI = "com.seonology.journey:/oauth2redirect"

    fun config(): AuthorizationServiceConfiguration {
        val issuer = BuildConfig.KEYCLOAK_ISSUER
        return AuthorizationServiceConfiguration(
            Uri.parse("$issuer/protocol/openid-connect/auth"),
            Uri.parse("$issuer/protocol/openid-connect/token"),
        )
    }

    fun buildAuthRequest(cfg: AuthorizationServiceConfiguration): AuthorizationRequest {
        return AuthorizationRequest.Builder(
            cfg,
            BuildConfig.KEYCLOAK_CLIENT_ID,
            ResponseTypeValues.CODE,
            Uri.parse(REDIRECT_URI),
        )
            .setScopes("openid", "profile", "email")
            .build()
    }

    fun handleTokenResponse(context: Context, resp: TokenResponse?) {
        val store = AuthStore(context)
        store.accessToken = resp?.accessToken
        store.refreshToken = resp?.refreshToken
    }

    fun newService(activity: Activity): AuthorizationService = AuthorizationService(activity)
}
