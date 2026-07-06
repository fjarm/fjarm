package xyz.fjarm.loginlibrary

import com.connectrpc.Code
import com.connectrpc.ConnectException
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AttemptLoginUseCaseImpl @Inject constructor(
    private val loginRepository: LoginRepository,
): AttemptLoginUseCase {

    companion object {
        private const val ERROR_GREETING = "RUH-ROH!"
        private const val HIDDEN_ERROR = "Something went wrong. Try again later."
        private const val UNAVAILABLE_MESSAGE = "Something went wrong on our end. Try again later."
    }

    override suspend fun invoke(email: String, password: String): Result<Unit> {
        // Use cases are responsible for being thread safe.
        // SEE: https://developer.android.com/kotlin/coroutines/coroutines-best-practices#main-safe
        return withContext(Dispatchers.IO) {
            try {
                // TODO(2026-07-06): Implement saving the session tokens to an encrypted DataStore.
                loginRepository.createSession(email, password)
                Result.success(Unit)
            } catch (e: ConnectException) {
                val message = when (e.code) {
                    Code.UNAVAILABLE -> UNAVAILABLE_MESSAGE
                    else -> {
                        // In a production app, consider logging this unexpected exception.
                        HIDDEN_ERROR
                    }
                }
                Result.failure(Throwable(message))
            } catch (_: Exception) {
                // Generic error handling is reserved for truly exceptional errors that happen
                // outside ConnectRPC.
                Result.failure(Throwable("$ERROR_GREETING $HIDDEN_ERROR"))
            }
        }
    }
}
