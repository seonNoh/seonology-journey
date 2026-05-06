package com.seonology.journey.data

import android.content.Context
import androidx.room.Room
import com.seonology.journey.BuildConfig
import com.seonology.journey.auth.AuthInterceptor
import com.seonology.journey.auth.AuthStore
import com.seonology.journey.data.local.JourneyDatabase
import com.squareup.moshi.Moshi
import com.squareup.moshi.kotlin.reflect.KotlinJsonAdapterFactory
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.moshi.MoshiConverterFactory

object Network {
    @Volatile private var dbInstance: JourneyDatabase? = null
    @Volatile private var apiInstance: JourneyApi? = null

    fun journeyApi(authStore: AuthStore): JourneyApi {
        apiInstance?.let { return it }
        synchronized(this) {
            apiInstance?.let { return it }
            val logging = HttpLoggingInterceptor().apply { level = HttpLoggingInterceptor.Level.BASIC }
            val client = OkHttpClient.Builder()
                .addInterceptor(AuthInterceptor(authStore))
                .addInterceptor(logging)
                .build()
            val moshi = Moshi.Builder().add(KotlinJsonAdapterFactory()).build()
            val retrofit = Retrofit.Builder()
                .baseUrl("${BuildConfig.API_BASE}/api/v1/")
                .client(client)
                .addConverterFactory(MoshiConverterFactory.create(moshi))
                .build()
            val created = retrofit.create(JourneyApi::class.java)
            apiInstance = created
            return created
        }
    }

    fun database(context: Context): JourneyDatabase {
        dbInstance?.let { return it }
        synchronized(this) {
            dbInstance?.let { return it }
            // fallbackToDestructiveMigration keeps the MVP honest: schema
            // bumps wipe local cache rather than forcing hand-written
            // migrations before any real user data exists.
            val created = Room.databaseBuilder(
                context.applicationContext,
                JourneyDatabase::class.java,
                "journey.db",
            ).fallbackToDestructiveMigration().build()
            dbInstance = created
            return created
        }
    }
}
