package app.phraselab

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import java.net.HttpURLConnection
import java.net.URL
import java.net.URLEncoder

@Serializable
data class Practice(val phraseId: String, val level: String, val tokens: List<String>, val prompt: String)
@Serializable
data class CheckRequest(val phraseId: String, val tokens: List<String>)
@Serializable
data class CheckResult(val correct: Boolean, val zh: String, val pinyin: String, val ja: String)

/** Android API adapter. Add INTERNET permission and kotlinx-serialization in the host app. */
class PhraseApi(private val baseUrl: String = "http://10.0.2.2:8080") {
    private val json = Json { ignoreUnknownKeys = true }

    suspend fun newPractice(level: String): Practice = withContext(Dispatchers.IO) {
        val encoded = URLEncoder.encode(level, Charsets.UTF_8)
        get("/api/practice/new?level=$encoded").let { json.decodeFromString<Practice>(it) }
    }

    suspend fun check(phraseId: String, tokens: List<String>): CheckResult = withContext(Dispatchers.IO) {
        val connection = URL(baseUrl + "/api/practice/check").openConnection() as HttpURLConnection
        connection.requestMethod = "POST"
        connection.connectTimeout = 3000
        connection.readTimeout = 3000
        connection.doOutput = true
        connection.setRequestProperty("Content-Type", "application/json; charset=utf-8")
        connection.outputStream.use { it.write(json.encodeToString(CheckRequest(phraseId, tokens)).toByteArray()) }
        val status = connection.responseCode
        val stream = if (status in 200..299) connection.inputStream else connection.errorStream
        val body = stream.bufferedReader().use { it.readText() }
        connection.disconnect()
        check(status in 200..299) { "Phrase API returned HTTP $status" }
        json.decodeFromString<CheckResult>(body)
    }

    private fun get(path: String): String {
        val connection = URL(baseUrl + path).openConnection() as HttpURLConnection
        connection.connectTimeout = 3000
        connection.readTimeout = 3000
        val status = connection.responseCode
        val stream = if (status in 200..299) connection.inputStream else connection.errorStream
        val body = stream.bufferedReader().use { it.readText() }
        connection.disconnect()
        check(status in 200..299) { "Phrase API returned HTTP $status" }
        return body
    }
}
