package xyz.fjarm.loginandsignupimpl

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.systemBarsPadding
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.core.content.ContextCompat.getString
import androidx.hilt.lifecycle.viewmodel.compose.hiltViewModel
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.compose.LocalLifecycleOwner
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.repeatOnLifecycle
import xyz.fjarm.buttons.FjarmButton
import xyz.fjarm.buttons.FjarmOutlinedButton
import xyz.fjarm.fjarmtheme.FjarmTheme
import xyz.fjarm.previews.PreviewFontScales
import xyz.fjarm.previews.PreviewLightDarkTheme
import xyz.fjarm.previews.PreviewPhoneSizes
import xyz.fjarm.text.FjarmHeaderText
import xyz.fjarm.text.FjarmSubtitleText

@Composable
fun LoginAndSignUpScreen(
    modifier: Modifier = Modifier,
    viewModel: LoginAndSignUpViewModel = hiltViewModel(),
    navigateToSignUp: () -> Unit = {},
    navigateToLogIn: () -> Unit = {},
) {
    val lifecycleOwner = LocalLifecycleOwner.current

    val sideEffects = viewModel.sideEffect
    LaunchedEffect(sideEffects) {
        lifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
            sideEffects.collect { sideEffect ->
                when (sideEffect) {
                    is LoginAndSignUpSideEffect.NavigateToSignUp -> {
                        navigateToSignUp()
                    }
                    is LoginAndSignUpSideEffect.NavigateToLogIn -> {
                        navigateToLogIn()
                    }
                }
            }
        }
    }

    val context = LocalContext.current
    val state by viewModel.state.collectAsStateWithLifecycle()

    LoginAndSignUpContent(
        modifier = modifier,
        titleLine = getString(context, state.titleLineText),
        subtitleLine = getString(context, state.subtitleLineText),
        logo = state.logo,
        signUpButtonText = getString(context, state.signUpButtonText),
        logInButtonText = getString(context, state.logInButtonText),
        onJoinClick = { viewModel.processEvent(LoginAndSignUpEvent.SignUpButtonClicked) },
        onLoginClick = { viewModel.processEvent(LoginAndSignUpEvent.LogInButtonClicked) }
    )
}

@Composable
private fun LoginAndSignUpContent(
    modifier: Modifier = Modifier,
    titleLine: String = "Fjarm",
    subtitleLine: String = "Plan and execute workouts",
    logo: Int = R.drawable.runner_illustration,
    signUpButtonText: String = "Join for free",
    logInButtonText: String = "Log in",
    onJoinClick: () -> Unit = {},
    onLoginClick: () -> Unit = {},
) {
    FjarmTheme {
        Surface(
            modifier = modifier
                .fillMaxSize(),
        ) {
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    // Use systemBarsPadding to ensure content stays within safe areas while the
                    // [Surface] background bleeds edge-to-edge
                    .systemBarsPadding()
                    .padding(horizontal = 32.dp, vertical = 16.dp),
                horizontalAlignment = Alignment.CenterHorizontally,
                // SpacedBy ensures the header and buttons have distance if the screen is small
                verticalArrangement = Arrangement.spacedBy(16.dp),
            ) {
                // 1. Header
                FjarmHeaderText(
                    text = titleLine,
                    modifier = Modifier
                        .fillMaxWidth(),
                )
                FjarmSubtitleText(
                    text = subtitleLine,
                    modifier = Modifier
                        .fillMaxWidth(),
                )

                // 2. Illustration (Takes up all available remaining space)
                Box(
                    modifier = Modifier
                        .size(280.dp)
                        .weight(1f),
                    contentAlignment = Alignment.Center,
                ) {
                    // Illustration goes here
                    Image(
                        painter = painterResource(id = logo),
                        contentDescription = null,
                        modifier = Modifier.fillMaxSize(),
                    )
                }

                // 3. Buttons
                FjarmOutlinedButton(
                    onClick = onJoinClick,
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(56.dp),
                ) {
                    Text(
                        text = signUpButtonText,
                        fontWeight = FontWeight.Bold,
                        fontSize = 16.sp,
                    )
                }

                FjarmButton(
                    onClick = onLoginClick,
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(56.dp),
                ) {
                    Text(
                        text = logInButtonText,
                        fontWeight = FontWeight.SemiBold,
                        fontSize = 16.sp,
                    )
                }
            }
        }
    }
}

@PreviewPhoneSizes
@PreviewFontScales
@PreviewLightDarkTheme
@Composable
fun LoginAndSignUpScreenPreview() {
    LoginAndSignUpContent()
}
