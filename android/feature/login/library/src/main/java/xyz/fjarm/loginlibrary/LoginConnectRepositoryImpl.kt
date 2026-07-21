package xyz.fjarm.loginlibrary

import build.buf.gen.fjarm.authentication.v1.AuthenticationServiceClientInterface
import build.buf.gen.fjarm.authentication.v1.CreateSessionRequest
import build.buf.gen.fjarm.authentication.v1.Session
import build.buf.gen.fjarm.users.v1.UserEmailAddress
import build.buf.gen.fjarm.users.v1.UserPassword
import com.connectrpc.getOrThrow
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
        val request = CreateSessionRequest.newBuilder()
            .setEmailAddress(UserEmailAddress.newBuilder().setEmailAddress(email).build())
            .setPassword(UserPassword.newBuilder().setPassword(password).build())
            .build()
        val response = client.createSession(request, headers = emptyMap()).getOrThrow()
        return response.session
    }
}
