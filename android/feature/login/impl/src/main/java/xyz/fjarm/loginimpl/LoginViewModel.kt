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
                emailInput = "",
                emailInputIsValid = true,
                passwordInputLabelText = R.string.password_input_label_text,
                passwordInput = "",
            ),
            loginButton = LoginState.LoginButton(
                loginButtonText = R.string.login_button_text,
                loginButtonEnabled = true,
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
                _state.update {
                    it.copy(
                        userInput = it.userInput.copy(
                            emailInput = event.emailAddress,
                        )
                    )
                }
            }
            is LoginEvent.LoginButtonClicked -> {
                // TODO
            }
            is LoginEvent.PasswordModified -> {
                _state.update {
                    it.copy(
                        userInput = it.userInput.copy(
                            passwordInput = event.password,
                        )
                    )
                }
            }
        }
    }
}
