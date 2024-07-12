package xyz.fjarm.libhelloworld

import build.buf.gen.helloworld.v1.GetHelloWorldResponse
import build.buf.gen.helloworld.v1.HelloWorldServiceGrpcKt

interface HelloWorldClient {

    suspend fun getHelloWorld(): GetHelloWorldResponse
}
