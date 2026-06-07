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
import java.util.regex.Pattern
import javax.inject.Inject

private val emailPattern = Pattern.compile(
    "[a-zA-Z0-9\\+\\.\\_\\%\\-\\+]{1,256}" +
            "\\@" +
            "[a-zA-Z0-9][a-zA-Z0-9\\-]{0,64}" +
            "(" +
            "\\." +
            "[a-zA-Z0-9][a-zA-Z0-9\\-]{0,25}" +
            ")+"
)

@HiltViewModel
class LoginViewModel @Inject constructor(
    private val attemptLoginUseCase: AttemptLoginUseCase,
): ViewModel() {

    private val _state = MutableStateFlow<LoginState>(value = LoginState())
    val state = _state.asStateFlow()

    private val _sideEffect = MutableSharedFlow<LoginSideEffect>()
    val sideEffect = _sideEffect.asSharedFlow()

    fun processEvent(event: LoginEvent) {
        val action = mapEventToAction(event)
        handleAction(action)
    }

    private fun mapEventToAction(event: LoginEvent): LoginAction {
        return when (event) {
            is LoginEvent.EmailAddressModified -> LoginAction.UpdateEmailAddress(
                email = event.emailAddress
            )
            is LoginEvent.PasswordModified -> LoginAction.UpdatePassword(password = event.password)
            is LoginEvent.LoginButtonClicked -> LoginAction.ExecuteLogin(
                email = _state.value.userInput.emailInputText,
                password = _state.value.userInput.passwordInputText,
            )
        }
    }

    private fun handleAction(action: LoginAction) {
        when (action) {
            is LoginAction.UpdateEmailAddress -> {
                reduce(mutation = LoginMutation.EmailUpdated(email = action.email))
            }
            is LoginAction.UpdatePassword -> {
                reduce(mutation = LoginMutation.PasswordUpdated(password = action.password))
            }
            is LoginAction.ExecuteLogin -> {
                viewModelScope.launch {
                    reduce(LoginMutation.Loading)
                    val result = attemptLoginUseCase(
                        email = action.email,
                        password = action.password,
                    )
                    result.fold(
                        onSuccess = {
                            reduce(LoginMutation.Success)
                            _sideEffect.emit(LoginSideEffect.NavigateToHome)
                        },
                        onFailure = { e ->
                            reduce(LoginMutation.Error)
                            _sideEffect.emit(
                                LoginSideEffect.ShowSnackbar(e.message ?: "Uh-oh. Try again.")
                            )
                        }
                    )
                }
            }
        }
    }

    private fun reduce(mutation: LoginMutation) {
        _state.update { currentState ->
            when (mutation) {
                is LoginMutation.EmailUpdated -> {
                    val email = mutation.email
                    val emailIsValid = emailPattern
                        .matcher(email)
                        .matches()
                    val password = currentState.userInput.passwordInputText

                    currentState.copy(
                        userInput = currentState.userInput.copy(
                            emailInputText = email,
                            emailInputIsInvalid = email.isNotEmpty() && !emailIsValid,
                        ),
                        loginButton = currentState.loginButton.copy(
                            loginButtonEnabled = password.isNotEmpty() && emailIsValid,
                        ),
                    )
                }
                is LoginMutation.PasswordUpdated -> {
                    val email = currentState.userInput.emailInputText
                    val emailIsValid = emailPattern
                        .matcher(email)
                        .matches()
                    val password = mutation.password

                    currentState.copy(
                        userInput = currentState.userInput.copy(
                            passwordInputText = password,
                        ),
                        loginButton = currentState.loginButton.copy(
                            loginButtonEnabled = password.isNotEmpty() && emailIsValid,
                        ),
                    )
                }
                is LoginMutation.Loading -> {
                    currentState.copy(
                        loadingIndicator = currentState.loadingIndicator.copy(
                            loadingIndicatorVisible = true,
                        ),
                        loginButton = currentState.loginButton.copy(loginButtonEnabled = false),
                    )
                }
                is LoginMutation.Success -> {
                    currentState.copy(
                        loadingIndicator = currentState.loadingIndicator.copy(
                            loadingIndicatorVisible = false,
                        ),
                    )
                }
                is LoginMutation.Error -> {
                    currentState.copy(
                        loadingIndicator = currentState.loadingIndicator.copy(
                            loadingIndicatorVisible = false,
                        ),
                        loginButton = currentState.loginButton.copy(loginButtonEnabled = true),
                    )
                }
            }
        }
    }
}
