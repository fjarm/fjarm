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

    companion object {
        private const val TITLE_LINE = "Fjarm"
        private const val SUBTITLE_LINE = "Plan and execute workouts"
        private const val LOGO = "android.resource://xyz.fjarm.loginandsignupfeat/fjarm_login_logo"
        private const val SIGN_UP_BUTTON_TEXT = "Join for free"
        private const val LOG_IN_BUTTON_TEXT = "Log in"
    }

    private val _state = MutableStateFlow<LoginAndSignUpState>(LoginAndSignUpState(
        titleLineText = TITLE_LINE,
        subtitleLineText = SUBTITLE_LINE,
        logo = LOGO,
        signUpButtonText = SIGN_UP_BUTTON_TEXT,
        logInButtonText = LOG_IN_BUTTON_TEXT,
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
