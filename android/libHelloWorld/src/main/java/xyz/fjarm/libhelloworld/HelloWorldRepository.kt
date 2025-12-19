package xyz.fjarm.libhelloworld

import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldRequest
import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldResponse

interface HelloWorldClient {

    suspend fun getHelloWorld(
        request: GetHelloWorldRequest,
    ): GetHelloWorldResponse
}
