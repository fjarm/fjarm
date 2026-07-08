package xyz.fjarm.loginlibrary

sealed class AttemptLoginException(cause: Throwable) : Exception(cause) {
    class InvalidCredentials(cause: Throwable) : AttemptLoginException(cause)
    class ServerUnavailable(cause: Throwable) : AttemptLoginException(cause)
}
