package xyz.fjarm.libservertransport

import com.connectrpc.ProtocolClientConfig
import com.connectrpc.SerializationStrategy
import com.connectrpc.extensions.GoogleJavaLiteProtobufStrategy
import com.connectrpc.protocols.NetworkProtocol

class ServerTransportModule {

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