package xyz.fjarm.loginlibrary

import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AttemptLoginUseCaseImpl @Inject constructor(): AttemptLoginUseCase {

    override suspend fun invoke(action: LoginAction.AttemptLoginWithCredentials) {
        TODO("Not yet implemented")
    }
}
