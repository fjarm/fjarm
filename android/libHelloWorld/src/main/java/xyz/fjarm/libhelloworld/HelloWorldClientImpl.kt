package xyz.fjarm.libhelloworld

import build.buf.gen.helloworld.v1.GetHelloWorldResponse
import build.buf.gen.helloworld.v1.HelloWorldServiceGrpcKt.HelloWorldServiceCoroutineStub
import com.google.protobuf.Empty
import io.grpc.ManagedChannel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.util.concurrent.TimeUnit

class HelloWorldClientImpl(
    channel: ManagedChannel,
) : HelloWorldClient {

    private val stub: HelloWorldServiceCoroutineStub = HelloWorldServiceCoroutineStub(channel)

    override suspend fun getHelloWorld(): GetHelloWorldResponse {
        return withContext(Dispatchers.IO) {
            val req = Empty.getDefaultInstance()
            stub.withDeadlineAfter(5, TimeUnit.SECONDS)
                .getHelloWorld(req)
        }
    }
}
