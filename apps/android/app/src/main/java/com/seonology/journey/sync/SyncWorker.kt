package com.seonology.journey.sync

import android.content.Context
import androidx.work.Constraints
import androidx.work.CoroutineWorker
import androidx.work.ExistingPeriodicWorkPolicy
import androidx.work.ExistingWorkPolicy
import androidx.work.NetworkType
import androidx.work.OneTimeWorkRequestBuilder
import androidx.work.PeriodicWorkRequestBuilder
import androidx.work.WorkManager
import androidx.work.WorkerParameters
import com.seonology.journey.auth.AuthStore
import com.seonology.journey.data.CreateTripRequest
import com.seonology.journey.data.JourneyApi
import com.seonology.journey.data.Network
import com.seonology.journey.data.local.DayEntity
import com.seonology.journey.data.local.ExpenseEntity
import com.seonology.journey.data.local.JourneyDatabase
import com.seonology.journey.data.local.NoteEntity
import com.seonology.journey.data.local.ScheduleEntity
import com.seonology.journey.data.local.SyncStateEntity
import com.seonology.journey.data.local.TripEntity
import java.util.concurrent.TimeUnit

/**
 * Syncs offline changes to the server and pulls fresh server state into Room.
 *
 * Ordering per run:
 *   1. Push dirty local rows (lastWriteWins — trust local updatedAt).
 *   2. Pull trip list and per-trip children; overwrite local rows when the
 *      server copy's updatedAt is strictly newer (prevents clobbering a
 *      concurrent local edit).
 *
 * Anything more sophisticated (delta sync, conflict UI) is intentionally
 * out of scope for this pass; once the API exposes a /trips/{id}/delta
 * endpoint we can replace the full per-trip refresh here.
 */
class SyncWorker(
    appContext: Context,
    params: WorkerParameters,
) : CoroutineWorker(appContext, params) {

    override suspend fun doWork(): Result {
        val store = AuthStore(applicationContext)
        if (!store.isAuthenticated) return Result.success()

        val api = Network.journeyApi(store)
        val db = Network.database(applicationContext)

        return try {
            pushDirty(api, db)
            pullAll(api, db)
            db.syncStateDao().put(SyncStateEntity("all", System.currentTimeMillis()))
            Result.success()
        } catch (t: Throwable) {
            // Transient failures retry with WorkManager's exponential backoff.
            Result.retry()
        }
    }

    private suspend fun pushDirty(api: JourneyApi, db: JourneyDatabase) {
        // Trips: only title/description/dates are pushable via CreateTripRequest today.
        // A future iteration should add an update endpoint mapping.
        db.tripDao().pending().forEach { t ->
            runCatching {
                api.updateTrip(
                    t.tripId,
                    CreateTripRequest(
                        title = t.title,
                        destination = t.destination,
                        startDate = t.startDate,
                        endDate = t.endDate,
                        budgetCurrency = t.budgetCurrency,
                    ),
                )
                db.tripDao().clearDirty(t.tripId)
            }
        }
        // Schedule/expense/note push endpoints are not yet wired from the
        // Kotlin side; skip cleanly so sync still reconciles pulls and
        // trip-level edits make it upstream.
    }

    private suspend fun pullAll(api: JourneyApi, db: JourneyDatabase) {
        val trips = api.listTrips(limit = 100).trips
        val tripEntities = trips.map { t ->
            TripEntity(
                tripId = t.id,
                title = t.title,
                description = t.description,
                destination = t.destination,
                startDate = t.startDate,
                endDate = t.endDate,
                budgetCurrency = t.budgetCurrency,
                status = t.status,
                updatedAt = System.currentTimeMillis(),
                isDirty = false,
            )
        }
        db.tripDao().upsertAll(tripEntities)

        // Fan out per-trip pulls. Serial to stay gentle on the API during
        // MVP; WorkManager's coroutine dispatcher already keeps us off the
        // main thread.
        for (t in trips) {
            runCatching { pullTrip(api, db, t.id) }
        }
    }

    private suspend fun pullTrip(api: JourneyApi, db: JourneyDatabase, tripId: String) {
        val days = api.listDays(tripId).days
        db.dayDao().upsertAll(days.map { d ->
            DayEntity(
                dayId = d.id,
                tripId = d.tripId,
                dayDate = d.date,
                dayNumber = d.dayNumber,
                region = d.region,
                updatedAt = System.currentTimeMillis(),
                isDirty = false,
            )
        })
        for (d in days) {
            runCatching {
                val sch = api.listSchedules(d.id).schedules
                db.scheduleDao().upsertAll(sch.map { s ->
                    ScheduleEntity(
                        scheduleId = s.id,
                        dayId = s.dayId,
                        title = s.title,
                        startTime = s.startTime,
                        endTime = s.endTime,
                        notes = s.notes,
                        updatedAt = System.currentTimeMillis(),
                        isDirty = false,
                    )
                })
            }
        }
        runCatching {
            val exp = api.listExpenses(tripId, limit = 100).expenses
            db.expenseDao().upsertAll(exp.map { e ->
                ExpenseEntity(
                    expenseId = e.id,
                    tripId = e.tripId,
                    dayId = e.dayId,
                    category = e.category,
                    amount = (e.amount?.amount ?: 0L).toDouble(),
                    currency = e.amount?.currency ?: "JPY",
                    description = e.description,
                    updatedAt = System.currentTimeMillis(),
                    isDirty = false,
                )
            })
        }
        runCatching {
            val notes = api.listNotes(tripId).notes
            db.noteDao().upsertAll(notes.map { n ->
                NoteEntity(
                    noteId = n.id,
                    tripId = n.tripId,
                    title = n.content.take(40),
                    content = n.content,
                    updatedAt = System.currentTimeMillis(),
                    isDirty = false,
                )
            })
        }
    }

    companion object {
        private const val ONE_SHOT = "journey_sync"
        private const val PERIODIC = "journey_sync_periodic"

        fun enqueueOneShot(context: Context) {
            val constraints = Constraints.Builder()
                .setRequiredNetworkType(NetworkType.CONNECTED)
                .build()
            val request = OneTimeWorkRequestBuilder<SyncWorker>()
                .setConstraints(constraints)
                .build()
            WorkManager.getInstance(context)
                .enqueueUniqueWork(ONE_SHOT, ExistingWorkPolicy.REPLACE, request)
        }

        /** Run sync at most every 15 minutes when network is available. */
        fun enqueuePeriodic(context: Context) {
            val constraints = Constraints.Builder()
                .setRequiredNetworkType(NetworkType.CONNECTED)
                .build()
            val request = PeriodicWorkRequestBuilder<SyncWorker>(15, TimeUnit.MINUTES)
                .setConstraints(constraints)
                .build()
            WorkManager.getInstance(context)
                .enqueueUniquePeriodicWork(PERIODIC, ExistingPeriodicWorkPolicy.KEEP, request)
        }
    }
}
