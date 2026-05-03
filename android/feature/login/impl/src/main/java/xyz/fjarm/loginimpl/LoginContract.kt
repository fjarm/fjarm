package xyz.fjarm.loginimpl

import androidx.annotation.StringRes

data class LoginState(
    val header: Header,
    val userInput: UserInput,
    val loginButton: LoginButton,
    val loadingIndicator: LoadingIndicator,
    val footer: Footer,
) {

    data class Header(
        @StringRes val headerText: Int,
    )

    data class UserInput(
        @StringRes val emailInputLabelText: Int,
        var emailInputText: String,
        val emailInputIsValid: Boolean,
        @StringRes val passwordInputLabelText: Int,
        var passwordInputText: String,
    )

    data class LoginButton(
        @StringRes val loginButtonText: Int,
        val loginButtonEnabled: Boolean,
    )

    data class LoadingIndicator(
        val loadingIndicatorVisible: Boolean,
    )

    data class Footer(
        @StringRes val alternativeOptionsText: Int,
        @StringRes val newToFjarmPromptText: Int,
        @StringRes val navigateToSignUpButtonText: Int,
        @StringRes val privacyPolicyText: Int,
        @StringRes val termsOfServiceText: Int,
    )
}

sealed interface LoginEvent {

    data class EmailAddressModified(
        val emailAddress: String,
    ): LoginEvent

    data class PasswordModified(
        val password: String,
    ): LoginEvent

    data object LoginButtonClicked: LoginEvent
}

sealed interface LoginSideEffect {

    data object NavigateToHome: LoginSideEffect

    data class ShowSnackbar(
        val message: String,
    ): LoginSideEffect
}
