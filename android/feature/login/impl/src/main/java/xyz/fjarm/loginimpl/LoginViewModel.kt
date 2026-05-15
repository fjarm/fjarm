package xyz.fjarm.loginimpl

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asSharedFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import xyz.fjarm.loginlibrary.AttemptLoginUseCase
import javax.inject.Inject

@HiltViewModel
class LoginViewModel @Inject constructor(
    private val attemptLoginUseCase: AttemptLoginUseCase,
): ViewModel() {

    private val _state = MutableStateFlow<LoginState>(value = LoginState())
    val state = _state.asStateFlow()

    private val _sideEffect = MutableSharedFlow<LoginSideEffect>()
    val sideEffect = _sideEffect.asSharedFlow()

    fun processEvent(event: LoginEvent) {
        when (event) {
            is LoginEvent.EmailAddressModified -> {
                val email = event.emailAddress
                val emailIsValid = android.util.Patterns.EMAIL_ADDRESS
                    .matcher(email)
                    .matches()
                val password = _state.value.userInput.passwordInputText

                _state.update {
                    it.copy(
                        userInput = it.userInput.copy(
                            emailInputText = email,
                            emailInputIsInvalid = email.isNotEmpty() && !emailIsValid,
                        ),
                        loginButton = it.loginButton.copy(
                            loginButtonEnabled = password.isNotEmpty() && emailIsValid,
                        ),
                    )
                }
            }
            is LoginEvent.LoginButtonClicked -> {
                _state.update {
                    it.copy(
                        loginButton = it.loginButton.copy(
                            loginButtonEnabled = false,
                        ),
                        loadingIndicator = it.loadingIndicator.copy(
                            loadingIndicatorVisible = true,
                        ),
                    )
                }

                val email = _state.value.userInput.emailInputText
                val password = _state.value.userInput.passwordInputText

                viewModelScope.launch {
                    attemptLoginUseCase(email = email, password = password)
                }
            }
            is LoginEvent.PasswordModified -> {
                val email = _state.value.userInput.emailInputText
                val emailIsValid = android.util.Patterns.EMAIL_ADDRESS
                    .matcher(email)
                    .matches()
                val password = event.password

                _state.update {
                    it.copy(
                        userInput = it.userInput.copy(
                            passwordInputText = password,
                        ),
                        loginButton = it.loginButton.copy(
                            loginButtonEnabled = password.isNotEmpty() && emailIsValid,
                        ),
                    )
                }
            }
        }
    }
}
