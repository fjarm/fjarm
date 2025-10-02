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

    override suspend fun getHelloWorld(
        request: GetHelloWorldRequest,
    ): GetHelloWorldResponse {
        return withContext(Dispatchers.IO) {
            stub.withDeadlineAfter(5, TimeUnit.SECONDS)
                .getHelloWorld(request)
        }
    }
}
