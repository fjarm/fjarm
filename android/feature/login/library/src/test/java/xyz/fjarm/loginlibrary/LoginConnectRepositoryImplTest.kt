package xyz.fjarm.loginlibrary

import build.buf.gen.fjarm.authentication.v1.AuthenticationServiceClientInterface
import build.buf.gen.fjarm.authentication.v1.CreateSessionRequest
import build.buf.gen.fjarm.authentication.v1.CreateSessionResponse
import build.buf.gen.fjarm.authentication.v1.createSessionResponse
import build.buf.gen.fjarm.authentication.v1.session
import com.connectrpc.Code
import com.connectrpc.ConnectException
import com.connectrpc.Headers
import com.connectrpc.ResponseMessage
import com.connectrpc.http.Cancelable
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.runTest
import org.junit.Test
import xyz.fjarm.servertransport.IDEMPOTENCY_KEY
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNotNull
import kotlin.test.assertTrue

@OptIn(ExperimentalCoroutinesApi::class)
class LoginConnectRepositoryImplTest {

    @Test
    fun givenValidCredentials_whenCreateSession_thenSetIdempotencyKeyInRequestAndHeaders() = runTest(UnconfinedTestDispatcher()) {
        // Given
        var capturedRequest: CreateSessionRequest? = null
        var capturedHeaders: Headers? = null

        val fakeClient = object : AuthenticationServiceClientInterface {
            override suspend fun createSession(
                request: CreateSessionRequest,
                headers: Headers,
            ): ResponseMessage<CreateSessionResponse> {
                capturedRequest = request
                capturedHeaders = headers
                return ResponseMessage.Success(
                    message = createSessionResponse {
                        session = session { }
                    },
                    headers = emptyMap(),
                    trailers = emptyMap(),
                )
            }

            override fun createSession(
                request: CreateSessionRequest,
                headers: Headers,
                onResult: (ResponseMessage<CreateSessionResponse>) -> Unit
            ): Cancelable {
                throw NotImplementedError()
            }
        }
        val repository = LoginConnectRepositoryImpl(client = fakeClient)

        // When
        repository.createSession("user@example.com", "secret123")

        // Then
        val req = requireNotNull(capturedRequest)
        val headers = requireNotNull(capturedHeaders)
        val requestIdempotencyKey = req.idempotencyKey
        assertTrue(requestIdempotencyKey.isNotEmpty())

        val headerValues = headers[IDEMPOTENCY_KEY]
        assertNotNull(headerValues)
        assertEquals(listOf(requestIdempotencyKey), headerValues)
    }

    @Test
    fun givenValidCredentials_whenCreateSession_thenPassEmailAndPasswordInRequest() = runTest(UnconfinedTestDispatcher()) {
        // Given
        var capturedRequest: CreateSessionRequest? = null

        val fakeClient = object : AuthenticationServiceClientInterface {
            override suspend fun createSession(
                request: CreateSessionRequest,
                headers: Headers,
            ): ResponseMessage<CreateSessionResponse> {
                capturedRequest = request
                return ResponseMessage.Success(
                    message = createSessionResponse {
                        session = session { }
                    },
                    headers = emptyMap(),
                    trailers = emptyMap(),
                )
            }

            override fun createSession(
                request: CreateSessionRequest,
                headers: Headers,
                onResult: (ResponseMessage<CreateSessionResponse>) -> Unit
            ): Cancelable {
                throw NotImplementedError()
            }
        }
        val repository = LoginConnectRepositoryImpl(client = fakeClient)

        // When
        repository.createSession("user@example.com", "secret123")

        // Then
        val req = requireNotNull(capturedRequest)
        assertEquals("user@example.com", req.emailAddress.emailAddress)
        assertEquals("secret123", req.password.password)
    }

    @Test
    fun givenSuccessfulResponse_whenCreateSession_thenReturnSession() = runTest(UnconfinedTestDispatcher()) {
        // Given
        val expectedSession = session { }
        val fakeClient = object : AuthenticationServiceClientInterface {
            override suspend fun createSession(
                request: CreateSessionRequest,
                headers: Headers,
            ): ResponseMessage<CreateSessionResponse> {
                return ResponseMessage.Success(
                    message = createSessionResponse {
                        session = expectedSession
                    },
                    headers = emptyMap(),
                    trailers = emptyMap(),
                )
            }

            override fun createSession(
                request: CreateSessionRequest,
                headers: Headers,
                onResult: (ResponseMessage<CreateSessionResponse>) -> Unit
            ): Cancelable {
                throw NotImplementedError()
            }
        }
        val repository = LoginConnectRepositoryImpl(client = fakeClient)

        // When
        val resultSession = repository.createSession("user@example.com", "secret123")

        // Then
        assertEquals(expectedSession, resultSession)
    }

    @Test
    fun givenRpcFailure_whenCreateSession_thenThrowConnectException() = runTest(UnconfinedTestDispatcher()) {
        // Given
        val fakeClient = object : AuthenticationServiceClientInterface {
            override suspend fun createSession(
                request: CreateSessionRequest,
                headers: Headers,
            ): ResponseMessage<CreateSessionResponse> {
                return ResponseMessage.Failure(
                    cause = ConnectException(Code.UNAUTHENTICATED, "invalid credentials"),
                    headers = emptyMap(),
                    trailers = emptyMap(),
                )
            }

            override fun createSession(
                request: CreateSessionRequest,
                headers: Headers,
                onResult: (ResponseMessage<CreateSessionResponse>) -> Unit
            ): Cancelable {
                throw NotImplementedError()
            }
        }
        val repository = LoginConnectRepositoryImpl(client = fakeClient)

        // Then
        assertFailsWith<ConnectException> {
            // When
            repository.createSession("user@example.com", "wrongpassword")
        }
    }

    @Test
    fun givenMultipleCalls_whenCreateSession_thenGenerateUniqueIdempotencyKeys() = runTest(UnconfinedTestDispatcher()) {
        // Given
        val capturedKeys = mutableListOf<String>()

        val fakeClient = object : AuthenticationServiceClientInterface {
            override suspend fun createSession(
                request: CreateSessionRequest,
                headers: Headers,
            ): ResponseMessage<CreateSessionResponse> {
                capturedKeys.add(request.idempotencyKey)
                return ResponseMessage.Success(
                    message = createSessionResponse {
                        session = session { }
                    },
                    headers = emptyMap(),
                    trailers = emptyMap(),
                )
            }

            override fun createSession(
                request: CreateSessionRequest,
                headers: Headers,
                onResult: (ResponseMessage<CreateSessionResponse>) -> Unit
            ): Cancelable {
                throw NotImplementedError()
            }
        }
        val repository = LoginConnectRepositoryImpl(client = fakeClient)

        // When
        repository.createSession("user@example.com", "secret123")
        repository.createSession("user@example.com", "secret123")

        // Then
        assertEquals(2, capturedKeys.size)
        val firstKey = capturedKeys[0]
        val secondKey = capturedKeys[1]
        assertTrue(firstKey != secondKey, "Idempotency keys must be unique per request")
    }

    @Test
    fun givenValidCredentials_whenCreateSession_thenContainOnlyIdempotencyHeader() = runTest(UnconfinedTestDispatcher()) {
        // Given
        var capturedHeaders: Headers? = null

        val fakeClient = object : AuthenticationServiceClientInterface {
            override suspend fun createSession(
                request: CreateSessionRequest,
                headers: Headers,
            ): ResponseMessage<CreateSessionResponse> {
                capturedHeaders = headers
                return ResponseMessage.Success(
                    message = createSessionResponse {
                        session = session { }
                    },
                    headers = emptyMap(),
                    trailers = emptyMap(),
                )
            }

            override fun createSession(
                request: CreateSessionRequest,
                headers: Headers,
                onResult: (ResponseMessage<CreateSessionResponse>) -> Unit
            ): Cancelable {
                throw NotImplementedError()
            }
        }
        val repository = LoginConnectRepositoryImpl(client = fakeClient)

        // When
        repository.createSession("user@example.com", "secret123")

        // Then
        val headers = requireNotNull(capturedHeaders)
        assertEquals(1, headers.size, "Should only pass 1 header")
        assertTrue(headers.containsKey(IDEMPOTENCY_KEY), "Headers map must contain idempotency_key")
    }

    @Test
    fun givenServerUnavailable_whenCreateSession_thenThrowConnectException() = runTest(UnconfinedTestDispatcher()) {
        // Given
        val fakeClient = object : AuthenticationServiceClientInterface {
            override suspend fun createSession(
                request: CreateSessionRequest,
                headers: Headers,
            ): ResponseMessage<CreateSessionResponse> {
                return ResponseMessage.Failure(
                    cause = ConnectException(Code.UNAVAILABLE, "server unavailable"),
                    headers = emptyMap(),
                    trailers = emptyMap(),
                )
            }

            override fun createSession(
                request: CreateSessionRequest,
                headers: Headers,
                onResult: (ResponseMessage<CreateSessionResponse>) -> Unit
            ): Cancelable {
                throw NotImplementedError()
            }
        }
        val repository = LoginConnectRepositoryImpl(client = fakeClient)

        // Then
        val exception = assertFailsWith<ConnectException> {
            // When
            repository.createSession("user@example.com", "secret123")
        }
        assertEquals(Code.UNAVAILABLE, exception.code)
    }
}
