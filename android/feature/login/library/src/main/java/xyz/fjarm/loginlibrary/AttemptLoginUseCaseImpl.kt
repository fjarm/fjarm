package xyz.fjarm.loginlibrary

import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AttemptLoginUseCaseImpl @Inject constructor(): AttemptLoginUseCase {

    override suspend fun invoke(email: String, password: String): Result<Unit> {
        // TODO(2026-05-15): Update this with real logic.
        return Result.success(Unit)
    }
}
