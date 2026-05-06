package com.seonology.journey.data

import com.squareup.moshi.JsonClass
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.PATCH
import retrofit2.http.POST
import retrofit2.http.Path
import retrofit2.http.Query

// --- DTOs ---

@JsonClass(generateAdapter = true)
data class Trip(
    val id: String = "",
    val title: String = "",
    val description: String? = null,
    val destination: String? = null,
    val startDate: String? = null,
    val endDate: String? = null,
    val budgetCurrency: String? = null,
    val status: String? = null,
    val updatedAt: String? = null,
)

@JsonClass(generateAdapter = true)
data class PageInfo(
    val nextCursor: String? = null,
    val hasMore: Boolean = false,
)

@JsonClass(generateAdapter = true)
data class ListTripsResponse(
    val trips: List<Trip> = emptyList(),
    val page: PageInfo? = null,
)

@JsonClass(generateAdapter = true)
data class CreateTripRequest(
    val title: String,
    val destination: String? = null,
    val startDate: String? = null,
    val endDate: String? = null,
    val budgetCurrency: String? = null,
)

@JsonClass(generateAdapter = true)
data class TripEnvelope(val trip: Trip)

@JsonClass(generateAdapter = true)
data class Day(
    val id: String = "",
    val tripId: String = "",
    val date: String = "",
    val dayNumber: Int = 0,
    val region: String? = null,
    val weather: String? = null,
    val dailySummary: String? = null,
    val updatedAt: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListDaysResponse(val days: List<Day> = emptyList())

@JsonClass(generateAdapter = true)
data class Schedule(
    val id: String = "",
    val dayId: String = "",
    val order: Int = 0,
    val title: String = "",
    val startTime: String? = null,
    val endTime: String? = null,
    val placeName: String? = null,
    val notes: String? = null,
    val updatedAt: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListSchedulesResponse(val schedules: List<Schedule> = emptyList())

@JsonClass(generateAdapter = true)
data class Money(val currency: String, val amount: Long)

@JsonClass(generateAdapter = true)
data class Expense(
    val id: String = "",
    val tripId: String = "",
    val dayId: String? = null,
    val category: String = "",
    val amount: Money? = null,
    val description: String? = null,
    val updatedAt: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListExpensesResponse(
    val expenses: List<Expense> = emptyList(),
    val page: PageInfo? = null,
)

@JsonClass(generateAdapter = true)
data class Note(
    val id: String = "",
    val tripId: String = "",
    val dayId: String? = null,
    val content: String = "",
    val mood: String? = null,
    val updatedAt: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListNotesResponse(val notes: List<Note> = emptyList())

@JsonClass(generateAdapter = true)
data class Media(
    val id: String = "",
    val tripId: String = "",
    val dayId: String? = null,
    val s3Key: String = "",
    val thumbnailS3Key: String? = null,
    val mimeType: String? = null,
    val takenAt: String? = null,
    val caption: String? = null,
)

@JsonClass(generateAdapter = true)
data class ListMediaResponse(
    val items: List<Media> = emptyList(),
    val page: PageInfo? = null,
)

// --- API ---
//
// Endpoints exposed here mirror the REST gateway in apps/api. Only the
// methods actually used by the Android client are listed; adding more is
// a matter of appending signatures, no plumbing change needed.

interface JourneyApi {
    // Trips
    @GET("trips")
    suspend fun listTrips(
        @Query("cursor") cursor: String? = null,
        @Query("limit") limit: Int? = null,
    ): ListTripsResponse

    @POST("trips")
    suspend fun createTrip(@Body req: CreateTripRequest): TripEnvelope

    @GET("trips/{tripId}")
    suspend fun getTrip(@Path("tripId") tripId: String): TripEnvelope

    @PATCH("trips/{tripId}")
    suspend fun updateTrip(@Path("tripId") tripId: String, @Body req: CreateTripRequest): TripEnvelope

    @DELETE("trips/{tripId}")
    suspend fun deleteTrip(@Path("tripId") tripId: String)

    // Days
    @GET("trips/{tripId}/days")
    suspend fun listDays(@Path("tripId") tripId: String): ListDaysResponse

    // Schedules
    @GET("days/{dayId}/schedules")
    suspend fun listSchedules(@Path("dayId") dayId: String): ListSchedulesResponse

    @POST("days/{dayId}/schedules")
    suspend fun createSchedule(@Path("dayId") dayId: String, @Body body: Schedule): Schedule

    // Expenses
    @GET("trips/{tripId}/expenses")
    suspend fun listExpenses(
        @Path("tripId") tripId: String,
        @Query("cursor") cursor: String? = null,
        @Query("limit") limit: Int? = null,
    ): ListExpensesResponse

    // Notes
    @GET("trips/{tripId}/notes")
    suspend fun listNotes(@Path("tripId") tripId: String): ListNotesResponse

    // Media
    @GET("trips/{tripId}/media")
    suspend fun listMedia(
        @Path("tripId") tripId: String,
        @Query("day_id") dayId: String? = null,
        @Query("cursor") cursor: String? = null,
        @Query("limit") limit: Int? = null,
    ): ListMediaResponse
}
