package xyz.fjarm.loginlibrary

import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent

interface AttemptLoginUseCase {

    suspend operator fun invoke(action: LoginAction.AttemptLoginWithCredentials)

    @Module
    @InstallIn(SingletonComponent::class)
    interface AttemptLoginUseCaseModule {

        fun bindAttemptLoginUseCase(impl: AttemptLoginUseCaseImpl): AttemptLoginUseCase
    }
}
