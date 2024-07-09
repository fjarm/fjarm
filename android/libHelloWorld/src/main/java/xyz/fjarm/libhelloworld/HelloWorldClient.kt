package xyz.fjarm.libhelloworld

import build.buf.gen.helloworld.v1.GetHelloWorldResponse

interface HelloWorldClient {

    suspend fun getHelloWorld(): GetHelloWorldResponse
}
