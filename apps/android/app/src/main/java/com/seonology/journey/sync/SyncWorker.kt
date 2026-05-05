package com.seonology.journey.sync

import android.content.Context
import androidx.work.Constraints
import androidx.work.CoroutineWorker
import androidx.work.ExistingWorkPolicy
import androidx.work.NetworkType
import androidx.work.OneTimeWorkRequestBuilder
import androidx.work.WorkManager
import androidx.work.WorkerParameters

/**
 * Sync offline changes to server when network is available.
 * Uses lastWriteWins (updatedAt) conflict resolution.
 */
class SyncWorker(
    appContext: Context,
    params: WorkerParameters,
) : CoroutineWorker(appContext, params) {

    override suspend fun doWork(): Result {
        // TODO: Implement actual sync logic:
        // 1. Read pending changes from Room (updatedAt > lastSyncedAt)
        // 2. Push to server API (DDB-backed)
        // 3. Pull remote changes
        // 4. Resolve conflicts with lastWriteWins
        return Result.success()
    }

    companion object {
        private const val WORK_NAME = "journey_sync"

        fun enqueue(context: Context) {
            val constraints = Constraints.Builder()
                .setRequiredNetworkType(NetworkType.CONNECTED)
                .build()

            val request = OneTimeWorkRequestBuilder<SyncWorker>()
                .setConstraints(constraints)
                .build()

            WorkManager.getInstance(context)
                .enqueueUniqueWork(WORK_NAME, ExistingWorkPolicy.REPLACE, request)
        }
    }
}
