package xyz.fjarm.loginlibrary

import build.buf.gen.fjarm.authentication.v1.Session
import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent

interface LoginRepository {

    suspend fun createSession(email: String, password: String): Session

    @Module
    @InstallIn(SingletonComponent::class)
    interface LoginRepositoryModule {

        @Binds
        fun bindLoginConnectRepository(impl: LoginConnectRepositoryImpl): LoginRepository
    }
}
