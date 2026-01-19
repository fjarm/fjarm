package xyz.fjarm.libhelloworld

import build.buf.gen.fjarm.helloworld.v1.HelloWorldOutput

interface GetHelloWorldUseCase {

    suspend operator fun invoke(): HelloWorldOutput
}
