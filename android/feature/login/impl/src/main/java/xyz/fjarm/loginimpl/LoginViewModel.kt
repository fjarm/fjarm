package xyz.fjarm.loginimpl

import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asSharedFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import javax.inject.Inject

@HiltViewModel
class LoginViewModel @Inject constructor(): ViewModel() {

    private val _state = MutableStateFlow<LoginState>(
        LoginState(
            header = LoginState.Header(headerText = R.string.header_text),
            userInput = LoginState.UserInput(
                emailInputLabelText = R.string.email_input_label_text,
                emailInputText = "",
                emailInputIsInvalid = false,
                passwordInputLabelText = R.string.password_input_label_text,
                passwordInputText = "",
            ),
            loginButton = LoginState.LoginButton(
                loginButtonText = R.string.login_button_text,
                loginButtonEnabled = false,
            ),
            loadingIndicator = LoginState.LoadingIndicator(loadingIndicatorVisible = false),
            footer = LoginState.Footer(
                alternativeOptionsText = R.string.alternative_options_section_header_text,
                newToFjarmPromptText = R.string.new_to_fjarm_prompt_text,
                navigateToSignUpButtonText = R.string.navigate_to_sign_up_button_text,
                privacyPolicyText = R.string.navigate_to_privacy_policy_text,
                termsOfServiceText = R.string.navigate_to_terms_of_service_text,
            ),
        )
    )
    val state = _state.asStateFlow()

    private val _sideEffect = MutableSharedFlow<LoginSideEffect>()
    val sideEffect = _sideEffect.asSharedFlow()

    fun processEvent(event: LoginEvent) {
        when (event) {
            is LoginEvent.EmailAddressModified -> {
                val newState = reduce(
                    LoginMutation.EmailAddressModified(event.emailAddress),
                    _state.value,
                )
                _state.update { newState }
            }
            is LoginEvent.LoginButtonClicked -> {
                val newState = reduce(
                    LoginMutation.LoginButtonClicked,
                    _state.value,
                )
                _state.update { newState }
            }
            is LoginEvent.PasswordModified -> {
                val newState = reduce(
                    LoginMutation.PasswordModified(event.password),
                    _state.value,
                )
                _state.update { newState }
            }
        }
    }

    private fun reduce(mutation: LoginMutation, oldState: LoginState): LoginState {
        when (mutation) {
            is LoginMutation.EmailAddressModified -> {
                val email = mutation.emailAddress
                val emailIsValid = android.util.Patterns.EMAIL_ADDRESS
                    .matcher(email)
                    .matches()
                val password = oldState.userInput.passwordInputText

                return oldState.copy(
                    userInput = oldState.userInput.copy(
                        emailInputText = email,
                        emailInputIsInvalid = email.isNotEmpty() && !emailIsValid,
                    ),
                    loginButton = oldState.loginButton.copy(
                        loginButtonEnabled = password.isNotEmpty() && emailIsValid,
                    ),
                )
            }
            is LoginMutation.PasswordModified -> {
                val email = oldState.userInput.emailInputText
                val emailIsValid = android.util.Patterns.EMAIL_ADDRESS
                    .matcher(email)
                    .matches()
                val password = mutation.password

                return oldState.copy(
                    userInput = oldState.userInput.copy(
                        passwordInputText = password,
                    ),
                    loginButton = oldState.loginButton.copy(
                        loginButtonEnabled = password.isNotEmpty() && emailIsValid,
                    ),
                )
            }
            is LoginMutation.LoginButtonClicked -> {
                return oldState.copy(
                    loginButton = oldState.loginButton.copy(
                        loginButtonEnabled = false,
                    ),
                    loadingIndicator = oldState.loadingIndicator.copy(
                        loadingIndicatorVisible = true,
                    ),
                )
            }
            is LoginMutation.LoginFailed -> {
                return oldState.copy(
                    loginButton = oldState.loginButton.copy(
                        loginButtonEnabled = true,
                    ),
                    loadingIndicator = oldState.loadingIndicator.copy(
                        loadingIndicatorVisible = false,
                    ),
                )
            }
            is LoginMutation.LoginSucceeded -> {
                return oldState.copy(
                    loadingIndicator = oldState.loadingIndicator.copy(
                        loadingIndicatorVisible = false,
                    ),
                )
            }
        }
    }
}
