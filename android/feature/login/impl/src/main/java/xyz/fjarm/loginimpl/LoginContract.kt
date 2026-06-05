package xyz.fjarm.loginimpl

import androidx.annotation.StringRes

data class LoginState(
    val header: Header = Header(headerText = R.string.header_text),
    val userInput: UserInput = UserInput(
        emailInputLabelText = R.string.email_input_label_text,
        emailInputText = "",
        emailInputIsInvalid = false,
        passwordInputLabelText = R.string.password_input_label_text,
        passwordInputText = "",
    ),
    val loginButton: LoginButton = LoginButton(
        loginButtonText = R.string.login_button_text,
        loginButtonEnabled = false,
    ),
    val loadingIndicator: LoadingIndicator = LoadingIndicator(loadingIndicatorVisible = false),
    val footer: Footer = Footer(
        alternativeOptionsText = R.string.alternative_options_section_header_text,
        newToFjarmPromptText = R.string.new_to_fjarm_prompt_text,
        navigateToSignUpButtonText = R.string.navigate_to_sign_up_button_text,
        privacyPolicyText = R.string.navigate_to_privacy_policy_text,
        termsOfServiceText = R.string.navigate_to_terms_of_service_text,
    ),
) {

    data class Header(@StringRes val headerText: Int)

    data class UserInput(
        @StringRes val emailInputLabelText: Int,
        var emailInputText: String,
        val emailInputIsInvalid: Boolean,
        @StringRes val passwordInputLabelText: Int,
        var passwordInputText: String,
    )

    data class LoginButton(
        @StringRes val loginButtonText: Int,
        val loginButtonEnabled: Boolean,
    )

    data class LoadingIndicator(val loadingIndicatorVisible: Boolean)

    data class Footer(
        @StringRes val alternativeOptionsText: Int,
        @StringRes val newToFjarmPromptText: Int,
        @StringRes val navigateToSignUpButtonText: Int,
        @StringRes val privacyPolicyText: Int,
        @StringRes val termsOfServiceText: Int,
    )
}

// What the system must execute (like background work)
sealed interface LoginAction {

    data class UpdateEmailAddress(val email: String): LoginAction
    data class UpdatePassword(val password: String): LoginAction
    data class ExecuteLogin(
        val email: String,
        val password: String,
    ): LoginAction
}

// What the user did (click a button, type some text, etc.)
sealed interface LoginEvent {

    data class EmailAddressModified(val emailAddress: String): LoginEvent
    data class PasswordModified(val password: String): LoginEvent
    data object LoginButtonClicked: LoginEvent
}

// How the state must change (outcome of an action)
sealed interface LoginMutation {

    data class EmailUpdated(val email: String): LoginMutation
    data class PasswordUpdated(val password: String): LoginMutation
    data object Loading: LoginMutation
    data object Success: LoginMutation
    data object Error: LoginMutation
}

sealed interface LoginSideEffect {

    data object NavigateToHome: LoginSideEffect
    data class ShowSnackbar(val message: String): LoginSideEffect
}
