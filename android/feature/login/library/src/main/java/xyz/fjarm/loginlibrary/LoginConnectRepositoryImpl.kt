package xyz.fjarm.loginlibrary

import build.buf.gen.fjarm.authentication.v1.AuthenticationServiceClientInterface
import build.buf.gen.fjarm.authentication.v1.CreateSessionRequest
import build.buf.gen.fjarm.authentication.v1.Session
import build.buf.gen.fjarm.users.v1.UserEmailAddress
import build.buf.gen.fjarm.users.v1.UserPassword
import com.connectrpc.getOrThrow
import java.util.UUID
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class LoginConnectRepositoryImpl @Inject constructor(
    private val client: AuthenticationServiceClientInterface,
): LoginRepository {

    override suspend fun createSession(
        email: String,
        password: String,
    ): Session {
        val idempotencyKey = UUID.randomUUID().toString()
        val request = CreateSessionRequest.newBuilder()
            .setIdempotencyKey(idempotencyKey)
            .setEmailAddress(UserEmailAddress.newBuilder().setEmailAddress(email).build())
            .setPassword(UserPassword.newBuilder().setPassword(password).build())
            .build()
        val response = client.createSession(
            request,
            headers = mapOf("idempotency_key" to listOf(idempotencyKey)),
        ).getOrThrow()
        return response.session
    }
}
