package xyz.fjarm.loginlibrary

sealed interface LoginAction {

    data class LoginWithCredentials(
        val emailAddress: String,
        val password: String,
    ): LoginAction
}
