package xyz.fjarm.libhelloworld

import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldRequest
import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldResponse
import build.buf.gen.fjarm.helloworld.v1.HelloWorldServiceGrpcKt
import io.grpc.ManagedChannel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.util.concurrent.TimeUnit

class HelloWorldGrpcClientImpl(
    channel: ManagedChannel,
) : HelloWorldClient {

    private val stub: HelloWorldServiceGrpcKt.HelloWorldServiceCoroutineStub =
        HelloWorldServiceGrpcKt.HelloWorldServiceCoroutineStub(channel)

    override suspend fun getHelloWorld(): GetHelloWorldResponse {
        return withContext(Dispatchers.IO) {
            val req = GetHelloWorldRequest.newBuilder().build()
            stub.withDeadlineAfter(5, TimeUnit.SECONDS)
                .getHelloWorld(req)
        }
    }
}
