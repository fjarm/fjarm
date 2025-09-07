package xyz.fjarm.libhelloworld

import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldResponse

interface HelloWorldClient {

    suspend fun getHelloWorld(): GetHelloWorldResponse
}
