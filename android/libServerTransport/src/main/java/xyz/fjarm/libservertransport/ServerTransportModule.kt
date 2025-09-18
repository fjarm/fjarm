package xyz.fjarm.libservertransport

import io.grpc.ManagedChannel
import io.grpc.ManagedChannelBuilder
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.asExecutor

class ServerTransportModule {

    fun provideGrpcServerTransport(
        serverAddress: String,
        serverPort: Int,
    ): ManagedChannel {
        return ManagedChannelBuilder
            .forAddress(
                serverAddress,
                serverPort
            )
            .executor(Dispatchers.IO.asExecutor())
            .build()
    }
}