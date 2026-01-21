package xyz.fjarm.libhelloworld

import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldRequest
import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldResponse
import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent

interface HelloWorldRepository {

    suspend fun getHelloWorld(
        request: GetHelloWorldRequest,
    ): GetHelloWorldResponse

    @Module
    @InstallIn(SingletonComponent::class)
    abstract class HelloWorldRepositoryModule {

        @Binds
        abstract fun bindHelloWorldRepository(
            helloWorldRepositoryImpl: HelloWorldConnectRepositoryImpl
        ): HelloWorldRepository
    }
}
