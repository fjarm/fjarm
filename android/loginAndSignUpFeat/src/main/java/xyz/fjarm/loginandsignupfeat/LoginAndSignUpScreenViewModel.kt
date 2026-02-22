package xyz.fjarm.loginandsignupfeat

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asSharedFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class LoginAndSignUpScreenViewModel @Inject constructor(
): ViewModel() {

    private val _state = MutableStateFlow<LoginAndSignUpState>(LoginAndSignUpState(
        titleLineText = R.string.title,
        subtitleLineText = R.string.subtitle,
        logo = R.drawable.runner_illustration,
        signUpButtonText = R.string.sign_up_button,
        logInButtonText = R.string.log_in_button,
    ))
    val state = _state.asStateFlow()

    private val _sideEffect = MutableSharedFlow<LoginAndSignUpSideEffect>()
    val sideEffect = _sideEffect.asSharedFlow()

    fun processEvent(event: LoginAndSignUpEvent) {
        when (event) {
            is LoginAndSignUpEvent.SignUpButtonClicked -> {
                viewModelScope.launch {
                    _sideEffect.emit(
                        LoginAndSignUpSideEffect.NavigateToSignUp
                    )
                }
            }
            is LoginAndSignUpEvent.LogInButtonClicked -> {
                viewModelScope.launch {
                    _sideEffect.emit(
                        LoginAndSignUpSideEffect.NavigateToLogIn
                    )
                }
            }
        }
    }
}
