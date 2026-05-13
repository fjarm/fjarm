package xyz.fjarm.loginlibrary

sealed interface LoginAction {

    data class AttemptLoginWithCredentials(
        val emailAddress: String,
        val password: String,
    ): LoginAction
}
