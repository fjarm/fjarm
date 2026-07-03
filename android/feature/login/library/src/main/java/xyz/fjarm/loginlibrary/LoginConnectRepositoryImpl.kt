package xyz.fjarm.loginlibrary

import build.buf.gen.fjarm.authentication.v1.AuthenticationServiceClientInterface
import build.buf.gen.fjarm.authentication.v1.CreateSessionRequest
import build.buf.gen.fjarm.authentication.v1.CreateSessionResponse
import build.buf.gen.fjarm.authentication.v1.Session
import build.buf.gen.fjarm.users.v1.UserEmailAddress
import build.buf.gen.fjarm.users.v1.UserPassword
import com.connectrpc.Headers
import com.connectrpc.Idempotency
import com.connectrpc.MethodSpec
import com.connectrpc.ProtocolClientInterface
import com.connectrpc.ResponseMessage
import com.connectrpc.StreamType
import com.connectrpc.getOrThrow
import com.connectrpc.http.Cancelable
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class LoginConnectRepositoryImpl @Inject constructor(
    private val client: ProtocolClientInterface,
): LoginRepository, AuthenticationServiceClientInterface {

    companion object {
        // No need to inject this path. Service and schema upgrades can be done with rolling deploys
        // on backend and force upgrade prompts on the app.
        private const val CREATE_SESSION_V1_PATH =
            "fjarm.authentication.v1.AuthenticationService/CreateSession"

        private const val CALLBACK_NOT_SUPPORTED_MESSAGE =
            "Connect client does not support callbacks"
    }

    override suspend fun createSession(
        email: String,
        password: String,
    ): Session {
        val request = CreateSessionRequest.newBuilder()
            .setEmailAddress(UserEmailAddress.newBuilder().setEmailAddress(email).build())
            .setPassword(UserPassword.newBuilder().setPassword(password).build())
            .build()
        val response = createSession(request, emptyMap()).getOrThrow()
        return response.session
    }

    override suspend fun createSession(
        request: CreateSessionRequest,
        headers: Headers,
    ): ResponseMessage<CreateSessionResponse> {
        return client.unary(
            request,
            headers,
            MethodSpec<CreateSessionRequest, CreateSessionResponse>(
                CREATE_SESSION_V1_PATH,
                CreateSessionRequest::class,
                CreateSessionResponse::class,
                StreamType.UNARY,
                Idempotency.IDEMPOTENT,
            )
        )
    }

    override fun createSession(
        request: CreateSessionRequest,
        headers: Headers,
        onResult: (ResponseMessage<CreateSessionResponse>) -> Unit,
    ): Cancelable {
        throw NotImplementedError(CALLBACK_NOT_SUPPORTED_MESSAGE)
    }
}
