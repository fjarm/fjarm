package xyz.fjarm.libhelloworld

import build.buf.gen.fjarm.helloworld.v1.HelloWorldOutput
import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent

interface GetHelloWorldUseCase {

    suspend operator fun invoke(): HelloWorldOutput

    @Module
    @InstallIn(SingletonComponent::class)
    abstract class GetHelloWorldUseCaseModule {

        @Binds
        abstract fun bindGetHelloWorldUseCase(
            getHelloWorldUseCaseImpl: GetHelloWorldUseCaseImpl
        ): GetHelloWorldUseCase
    }
}
