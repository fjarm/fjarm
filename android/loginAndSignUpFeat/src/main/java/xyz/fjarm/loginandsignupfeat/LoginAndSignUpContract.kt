package xyz.fjarm.loginandsignupfeat

import androidx.annotation.DrawableRes
import androidx.annotation.StringRes

data class LoginAndSignUpState(
    @StringRes val titleLineText: Int,
    @StringRes val subtitleLineText: Int,
    @DrawableRes val logo: Int,
    @StringRes val signUpButtonText: Int,
    @StringRes val logInButtonText: Int,
)

sealed class LoginAndSignUpEvent {

    data object SignUpButtonClicked : LoginAndSignUpEvent()

    data object LogInButtonClicked : LoginAndSignUpEvent()
}

sealed class LoginAndSignUpSideEffect {

    data object NavigateToSignUp : LoginAndSignUpSideEffect()

    data object NavigateToLogIn : LoginAndSignUpSideEffect()
}
