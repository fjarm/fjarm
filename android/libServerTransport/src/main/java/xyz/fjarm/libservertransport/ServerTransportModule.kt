package xyz.fjarm.libservertransport

import com.connectrpc.ProtocolClientConfig
import com.connectrpc.ProtocolClientInterface
import com.connectrpc.SerializationStrategy
import com.connectrpc.extensions.GoogleJavaLiteProtobufStrategy
import com.connectrpc.http.HTTPClientInterface
import com.connectrpc.impl.ProtocolClient
import com.connectrpc.okhttp.ConnectOkHttpClient
import com.connectrpc.protocols.NetworkProtocol
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import kotlinx.coroutines.Dispatchers
import okhttp3.OkHttpClient
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
class ServerTransportModule {

    @Provides
    @Singleton
    fun provideProtocolClient(
        httpClient: HTTPClientInterface,
        protocolClientConfig: ProtocolClientConfig,
    ): ProtocolClientInterface {
        return ProtocolClient(
            httpClient = httpClient,
            config = protocolClientConfig,
        )
    }

    @Provides
    @Singleton
    fun provideHttpClientInterface(
        okHttpClient: OkHttpClient,
    ): HTTPClientInterface {
        return ConnectOkHttpClient(
            unaryClient = okHttpClient
        )
    }

    @Provides
    @Singleton
    fun provideOkHttpClient(): OkHttpClient {
        // TODO(2025-09-23): Investigate setting up options like Authenticator, x509TrustManager, and/or Interceptor
        return OkHttpClient
            .Builder()
            .build()
    }

    @Provides
    @Singleton
    fun provideProtocolClientConfig(
        @ServerHost host: String,
        networkProtocol: NetworkProtocol,
        serializationStrategy: SerializationStrategy,
    ): ProtocolClientConfig {
        // TODO(2025-09-23): Investigate using request compression tools and adding interceptors.
        return ProtocolClientConfig(
            host = host,
            networkProtocol = networkProtocol,
            serializationStrategy = serializationStrategy,
            ioCoroutineContext = Dispatchers.IO,
        )
    }

    @Provides
    @Singleton
    @ServerHost
    fun provideHost(): String {
        // TODO(2025-09-23): Dagger inject an address instead of hardcoding localhost
        return "10.0.2.2"
    }

    @Provides
    @Singleton
    fun provideNetworkProtocol(): NetworkProtocol {
        return NetworkProtocol.CONNECT
    }

    @Provides
    @Singleton
    fun provideSerializationStrategy(): SerializationStrategy {
        return GoogleJavaLiteProtobufStrategy()
    }
}
