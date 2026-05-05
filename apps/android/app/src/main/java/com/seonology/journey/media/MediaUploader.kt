package com.seonology.journey.media

import android.content.Context
import android.net.Uri
import androidx.exifinterface.media.ExifInterface
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody
import okio.BufferedSink
import okio.source
import java.io.InputStream

data class ExifData(
    val lat: Double?,
    val lng: Double?,
    val takenAt: String?,
)

/** Upload file to S3 presigned URL with progress callback. */
suspend fun uploadToS3(
    presignedUrl: String,
    contentType: String,
    inputStream: InputStream,
    contentLength: Long,
    onProgress: (fraction: Float) -> Unit = {},
) = withContext(Dispatchers.IO) {
    val body = object : RequestBody() {
        override fun contentType() = contentType.toMediaType()
        override fun contentLength() = contentLength
        override fun writeTo(sink: BufferedSink) {
            val source = inputStream.source()
            var totalWritten = 0L
            val buffer = okio.Buffer()
            while (true) {
                val read = source.read(buffer, 8192)
                if (read == -1L) break
                sink.write(buffer, read)
                totalWritten += read
                onProgress(totalWritten.toFloat() / contentLength)
            }
        }
    }
    val request = Request.Builder().url(presignedUrl).put(body).build()
    OkHttpClient().newCall(request).execute().use { resp ->
        require(resp.isSuccessful) { "S3 upload failed: ${resp.code}" }
    }
}

/** Extract EXIF GPS and DateTime from image URI. */
fun extractExif(context: Context, uri: Uri): ExifData {
    val stream = context.contentResolver.openInputStream(uri) ?: return ExifData(null, null, null)
    return stream.use {
        val exif = ExifInterface(it)
        val latLong = exif.latLong
        ExifData(
            lat = latLong?.get(0),
            lng = latLong?.get(1),
            takenAt = exif.getAttribute(ExifInterface.TAG_DATETIME_ORIGINAL),
        )
    }
}
