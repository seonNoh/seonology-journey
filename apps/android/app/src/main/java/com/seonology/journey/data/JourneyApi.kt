package com.seonology.journey.data

import com.squareup.moshi.JsonClass
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.Path

@JsonClass(generateAdapter = true)
data class Trip(
    val tripId: String = "",
    val title: String = "",
    val description: String? = null,
    val destination: String? = null,
    val startDate: String? = null,
    val endDate: String? = null,
    val budgetCurrency: String? = null,
    val status: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListTripsResponse(val trips: List<Trip> = emptyList())

@JsonClass(generateAdapter = true)
data class CreateTripRequest(
    val title: String,
    val destination: String? = null,
    val startDate: String? = null,
    val endDate: String? = null,
    val budgetCurrency: String? = null,
)

@JsonClass(generateAdapter = true)
data class GetTripResponse(val trip: Trip)

@JsonClass(generateAdapter = true)
data class Day(
    val dayId: String,
    val tripId: String,
    val dayDate: String,
    val dayNumber: Int = 0,
    val region: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListDaysResponse(val days: List<Day> = emptyList())

interface JourneyApi {
    @GET("trips")
    suspend fun listTrips(): ListTripsResponse

    @POST("trips")
    suspend fun createTrip(@Body req: CreateTripRequest): GetTripResponse

    @GET("trips/{tripId}")
    suspend fun getTrip(@Path("tripId") tripId: String): GetTripResponse

    @GET("trips/{tripId}/days")
    suspend fun listDays(@Path("tripId") tripId: String): ListDaysResponse
}
