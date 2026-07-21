package xyz.fjarm.loginlibrary

import com.connectrpc.Code
import com.connectrpc.ConnectException
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AttemptLoginUseCaseImpl @Inject constructor(
    private val loginRepository: LoginRepository,
): AttemptLoginUseCase {

    override suspend fun invoke(email: String, password: String): Result<Unit> {
        // Use cases are responsible for being thread safe.
        // SEE: https://developer.android.com/kotlin/coroutines/coroutines-best-practices#main-safe
        return withContext(Dispatchers.IO) {
            try {
                // TODO(2026-07-06): Implement saving the session tokens to an encrypted DataStore.
                loginRepository.createSession(email, password)
                Result.success(Unit)
            } catch (e: ConnectException) {
                val failure = when (e.code) {
                    Code.UNAUTHENTICATED -> AttemptLoginException.InvalidCredentials(e)
                    Code.UNAVAILABLE -> AttemptLoginException.ServerUnavailable(e)
                    else -> e
                }
                Result.failure(failure)
            } catch (e: CancellationException) {
                throw e
            } catch (e: Exception) {
                // Generic error handling is reserved for truly exceptional errors that happen
                // outside ConnectRPC.
                Result.failure(e)
            }
        }
    }
}
