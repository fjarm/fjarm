package xyz.fjarm.libhelloworld

import build.buf.gen.helloworld.v1.GetHelloWorldResponse
import build.buf.gen.helloworld.v1.HelloWorldServiceGrpcKt.HelloWorldServiceCoroutineStub
import com.google.protobuf.Empty

class HelloWorldClientImpl(
    private val stub: HelloWorldServiceCoroutineStub,
) : HelloWorldClient {

    override suspend fun getHelloWorld(): GetHelloWorldResponse {
        val req = Empty.getDefaultInstance()
        return stub.getHelloWorld(req)
    }
}
