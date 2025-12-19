package xyz.fjarm.libhelloworld

import build.buf.gen.fjarm.helloworld.v1.HelloWorldOutput

interface HelloWorldInteractor {

    suspend fun getHelloWorld(): HelloWorldOutput
}
