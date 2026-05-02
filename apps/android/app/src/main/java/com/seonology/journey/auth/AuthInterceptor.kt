package com.seonology.journey.auth

import okhttp3.Interceptor
import okhttp3.Response

class AuthInterceptor(private val store: AuthStore) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val token = store.accessToken
        val req = if (!token.isNullOrEmpty())
            chain.request().newBuilder().header("Authorization", "Bearer $token").build()
        else chain.request()
        return chain.proceed(req)
    }
}
