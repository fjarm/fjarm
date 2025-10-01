package xyz.fjarm.libservertransport

import com.connectrpc.ProtocolClientConfig
import com.connectrpc.SerializationStrategy
import com.connectrpc.extensions.GoogleJavaLiteProtobufStrategy
import com.connectrpc.http.HTTPClientInterface
import com.connectrpc.impl.ProtocolClient
import com.connectrpc.okhttp.ConnectOkHttpClient
import com.connectrpc.protocols.NetworkProtocol
import kotlinx.coroutines.Dispatchers
import okhttp3.OkHttpClient

class ServerTransportModule {

    fun provideProtocolClient(
        httpClient: HTTPClientInterface,
        protocolClientConfig: ProtocolClientConfig,
    ): ProtocolClient {
        return ProtocolClient(
            httpClient = httpClient,
            config = protocolClientConfig,
        )
    }

    fun provideHttpClient(): HTTPClientInterface {
        return ConnectOkHttpClient(
            // TODO(2025-09-23): Investigate setting up options like Authenticator, x509TrustManager, and/or Interceptor
            unaryClient = OkHttpClient
                .Builder()
                .build()
        )
    }

    fun provideProtocolClientConfig(
        host: String,
        networkProtocol: NetworkProtocol,
        serializationStrategy: SerializationStrategy,
    ): ProtocolClientConfig {
        // TODO(2025-09-23): Investigate using the request compression tools.
        return ProtocolClientConfig(
            host = host,
            networkProtocol = networkProtocol,
            serializationStrategy = serializationStrategy,
            ioCoroutineContext = Dispatchers.IO,
        )
    }

    fun provideHost(): String {
        // TODO(2025-09-23): Dagger inject an address instead of hardcoding localhost
        return "10.0.2.2"
    }

    fun provideNetworkProtocol(): NetworkProtocol {
        return NetworkProtocol.CONNECT
    }

    fun provideSerializationStrategy(): SerializationStrategy {
        return GoogleJavaLiteProtobufStrategy()
    }
}