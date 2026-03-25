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
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
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
import xyz.fjarm.fjarmtheme.FjarmTheme
import xyz.fjarm.fjarmtheme.FjarmTypography
import xyz.fjarm.previews.PreviewFontScales
import xyz.fjarm.previews.PreviewLightDarkTheme
import xyz.fjarm.previews.PreviewPhoneSizes

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
                Text(
                    text = titleLine,
                    modifier = Modifier
                        .fillMaxWidth(),
                    style = FjarmTypography.headlineLarge.copy(
                        fontWeight = FontWeight.Bold,
                    ),
                )
                Text(
                    text = subtitleLine,
                    modifier = Modifier
                        .fillMaxWidth(),
                    style = FjarmTypography.bodyLarge,
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
                OutlinedButton(
                    onClick = onJoinClick,
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(56.dp),
                    shape = RoundedCornerShape(4.dp),
                    border = ButtonDefaults.outlinedButtonBorder(true).copy(width = 1.dp),
                ) {
                    Text(
                        text = signUpButtonText,
                        fontWeight = FontWeight.Bold,
                        fontSize = 16.sp,
                    )
                }

                TextButton(
                    onClick = onLoginClick,
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(56.dp),
                    shape = RoundedCornerShape(4.dp),
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
