package xyz.fjarm.loginlibrary

import build.buf.gen.fjarm.authentication.v1.Session
import build.buf.gen.fjarm.authentication.v1.session
import com.connectrpc.Code
import com.connectrpc.ConnectException
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.runTest
import org.junit.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertIs
import kotlin.test.assertTrue

@OptIn(ExperimentalCoroutinesApi::class)
class AttemptLoginUseCaseImplTest {

    @Test
    fun givenSuccessfulResponse_whenAttemptingLogin_thenReturnSuccess() = runTest(UnconfinedTestDispatcher()) {
        // Given
        val repository = object : LoginRepository {
            override suspend fun createSession(email: String, password: String): Session {
                return session { }
            }
        }
        val useCase = AttemptLoginUseCaseImpl(
            loginRepository = repository,
            ioDispatcher = UnconfinedTestDispatcher(testScheduler)
        )

        // When
        val result = useCase("email", "password")

        // Then
        assertTrue(result.isSuccess)
    }

    @Test
    fun givenInvalidCredentials_whenAttemptingLogin_thenReturnInvalidCredentialsFailure() = runTest(UnconfinedTestDispatcher()) {
        // Given
        val repository = object : LoginRepository {
            override suspend fun createSession(email: String, password: String): Session {
                throw ConnectException(Code.UNAUTHENTICATED)
            }
        }
        val useCase = AttemptLoginUseCaseImpl(
            loginRepository = repository,
            ioDispatcher = UnconfinedTestDispatcher(testScheduler)
        )

        // When
        val result = useCase("email", "password")

        // Then
        assertTrue(result.isFailure)
        assertIs<AttemptLoginException.InvalidCredentials>(result.exceptionOrNull())
    }

    @Test
    fun givenServerUnavailable_whenAttemptingLogin_thenReturnServerUnavailableFailure() = runTest(UnconfinedTestDispatcher()) {
        // Given
        val repository = object : LoginRepository {
            override suspend fun createSession(email: String, password: String): Session {
                throw ConnectException(Code.UNAVAILABLE)
            }
        }
        val useCase = AttemptLoginUseCaseImpl(
            loginRepository = repository,
            ioDispatcher = UnconfinedTestDispatcher(testScheduler)
        )

        // When
        val result = useCase("email", "password")

        // Then
        assertTrue(result.isFailure)
        assertIs<AttemptLoginException.ServerUnavailable>(result.exceptionOrNull())
    }

    @Test
    fun givenOtherConnectError_whenAttemptingLogin_thenReturnRawConnectException() = runTest(UnconfinedTestDispatcher()) {
        // Given
        val repository = object : LoginRepository {
            override suspend fun createSession(email: String, password: String): Session {
                throw ConnectException(Code.INTERNAL_ERROR)
            }
        }
        val useCase = AttemptLoginUseCaseImpl(
            loginRepository = repository,
            ioDispatcher = UnconfinedTestDispatcher(testScheduler)
        )

        // When
        val result = useCase("email", "password")

        // Then
        assertTrue(result.isFailure)
        val exception = result.exceptionOrNull()
        assertIs<ConnectException>(exception)
        assertEquals(Code.INTERNAL_ERROR, exception.code)
    }

    @Test
    fun givenCancellation_whenAttemptingLogin_thenRethrowCancellationException() = runTest(UnconfinedTestDispatcher()) {
        // Given
        val repository = object : LoginRepository {
            override suspend fun createSession(email: String, password: String): Session {
                throw CancellationException()
            }
        }
        val useCase = AttemptLoginUseCaseImpl(
            loginRepository = repository,
            ioDispatcher = UnconfinedTestDispatcher(testScheduler)
        )

        // Then
        assertFailsWith<CancellationException> {
            // When
            useCase("email", "password")
        }
    }

    @Test
    fun givenGenericException_whenAttemptingLogin_thenReturnGenericFailure() = runTest(UnconfinedTestDispatcher()) {
        // Given
        val repository = object : LoginRepository {
            override suspend fun createSession(email: String, password: String): Session {
                throw RuntimeException("Something went wrong")
            }
        }
        val useCase = AttemptLoginUseCaseImpl(
            loginRepository = repository,
            ioDispatcher = UnconfinedTestDispatcher(testScheduler)
        )

        // When
        val result = useCase("email", "password")

        // Then
        assertTrue(result.isFailure)
        assertIs<RuntimeException>(result.exceptionOrNull())
        assertEquals("Something went wrong", result.exceptionOrNull()?.message)
    }
}
