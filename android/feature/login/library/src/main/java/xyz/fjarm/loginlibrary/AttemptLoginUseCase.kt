package xyz.fjarm.loginlibrary

import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent

interface AttemptLoginUseCase {

    suspend operator fun invoke(email: String, password: String): Result<Unit>

    @Module
    @InstallIn(SingletonComponent::class)
    interface AttemptLoginUseCaseModule {

        fun bindAttemptLoginUseCase(impl: AttemptLoginUseCaseImpl): AttemptLoginUseCase
    }
}
