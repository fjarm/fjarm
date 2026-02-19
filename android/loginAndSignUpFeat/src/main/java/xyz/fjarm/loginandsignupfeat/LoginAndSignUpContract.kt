package xyz.fjarm.loginandsignupfeat

data class LoginAndSignUpState(
    val titleLine: String,
    val subtitleLine: String,
    // Construct a URI representing the drawable resource: android.resource://<package>/R.drawable.<drawable>
    val logo: String,
    val signUpButtonText: String,
    val logInButtonText: String,
)

sealed class LoginAndSignUpEvent {

    data object SignUpButtonClicked : LoginAndSignUpEvent()

    data object LogInButtonClicked : LoginAndSignUpEvent()
}

sealed class LoginAndSignUpSideEffect {

    data object NavigateToSignUp : LoginAndSignUpSideEffect()

    data object NavigateToLogIn : LoginAndSignUpSideEffect()
}
