package xyz.fjarm.loginlibrary

import build.buf.gen.fjarm.authentication.v1.AuthenticationServiceClient
import build.buf.gen.fjarm.authentication.v1.AuthenticationServiceClientInterface
import com.connectrpc.ProtocolClientInterface
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent

@Module
@InstallIn(SingletonComponent::class)
object AuthenticationServiceClientModule {

    @Provides
    fun provideAuthenticationServiceClient(
        client: ProtocolClientInterface,
    ): AuthenticationServiceClientInterface {
        return AuthenticationServiceClient(client)
    }
}
