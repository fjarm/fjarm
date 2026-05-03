package xyz.fjarm.loginimpl

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.defaultMinSize
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.layout.wrapContentHeight
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarDuration
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.SnackbarHostState
import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.tooling.preview.PreviewFontScale
import androidx.compose.ui.tooling.preview.PreviewLightDark
import androidx.compose.ui.tooling.preview.PreviewScreenSizes
import androidx.compose.ui.unit.dp
import androidx.hilt.lifecycle.viewmodel.compose.hiltViewModel
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.compose.LocalLifecycleOwner
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.repeatOnLifecycle
import xyz.fjarm.buttons.FjarmButton
import xyz.fjarm.buttons.FjarmFilledButton
import xyz.fjarm.fjarmtheme.FjarmTheme
import xyz.fjarm.text.FjarmExtraLargeHeaderText
import xyz.fjarm.text.FjarmNormalSizeText
import xyz.fjarm.text.FjarmSmallSizeText

@Composable
fun LoginScreen(
    modifier: Modifier = Modifier,
    viewModel: LoginViewModel = hiltViewModel(),
    navigateToHome: () -> Unit = {},
) {
    val lifecycleOwner = LocalLifecycleOwner.current
    val snackbarHostState = remember { SnackbarHostState() }

    val sideEffects = viewModel.sideEffect
    LaunchedEffect(sideEffects) {
        lifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
            sideEffects.collect { sideEffect ->
                when (sideEffect) {
                    is LoginSideEffect.NavigateToHome -> {
                        navigateToHome()
                    }
                    is LoginSideEffect.ShowSnackbar -> {
                        snackbarHostState.showSnackbar(
                            message = sideEffect.message,
                            duration = SnackbarDuration.Indefinite,
                        )
                    }
                }
            }
        }
    }

    val state by viewModel.state.collectAsStateWithLifecycle()

    Scaffold(
        containerColor = Color.Transparent,
        snackbarHost = {
            SnackbarHost(hostState = snackbarHostState)
        }
    ) { contentPadding ->
        LoginScreenContent(
            modifier = modifier.padding(contentPadding),
            emailInputText = state.userInput.emailInputText,
            onEmailInputTextModified = {
                viewModel.processEvent(LoginEvent.EmailAddressModified(it))
            },
            passwordInputText = state.userInput.passwordInputText,
            onPasswordInputTextModified = {
                viewModel.processEvent(LoginEvent.PasswordModified(it))
            },
        )
    }
}

@Composable
private fun LoginScreenContent(
    modifier: Modifier = Modifier,
    headerText: String = "Login",
    emailInputText: String = "",
    emailInputLabelText: String = "Email address",
    passwordInputText: String = "",
    passwordInputLabelText: String = "Password",
    loginButtonText: String = "Continue",
    loginButtonEnabled: Boolean = true,
    alternativeOptionsText: String = "Or",
    newToFjarmPromptText: String = "New to Fjarm?",
    navigateToSignUpButtonText: String = "Sign up",
    privacyPolicyText: String = "Privacy policy",
    termsOfServiceText: String = "Terms of service",
    onEmailInputTextModified: (String) -> Unit = {},
    onPasswordInputTextModified: (String) -> Unit = {},
) {
    val scrollState = rememberScrollState()

    FjarmTheme {
        Surface(
            modifier = modifier
                .fillMaxSize(),
        ) {
            Column(
                // No need to apply systemBarsPadding as the wrapping Scaffold already includes it.
                modifier = Modifier
                    .fillMaxSize()
                    .verticalScroll(scrollState)
                    .padding(horizontal = 16.dp, vertical = 16.dp),
                horizontalAlignment = Alignment.CenterHorizontally,
            ) {
                // Large Header
                Column(
                    modifier = Modifier
                        .fillMaxWidth(),
                ) {
                    FjarmExtraLargeHeaderText(
                        text = headerText,
                    )
                    Box(
                        modifier = Modifier
                            .width(124.dp)
                            .height(4.dp)
                            .background(MaterialTheme.colorScheme.onBackground),
                    )
                }

                Spacer(modifier = Modifier.height(64.dp))

                // Input Field Section
                OutlinedTextField(
                    value = emailInputText,
                    onValueChange = { onEmailInputTextModified(it) },
                    modifier = Modifier
                        .fillMaxWidth()
                        .defaultMinSize(minHeight = 56.dp),
                    label = {
                        FjarmNormalSizeText(
                            text = emailInputLabelText,
                            fontWeight = FontWeight.Bold,
                        )
                    },
                    shape = RoundedCornerShape(4.dp),
                    keyboardOptions = KeyboardOptions(
                        keyboardType = KeyboardType.Email,
                        imeAction = ImeAction.Next,
                    ),
                    singleLine = true,
                )

                OutlinedTextField(
                    value = passwordInputText,
                    onValueChange = {
                        onPasswordInputTextModified(it)
                    },
                    modifier = Modifier
                        .fillMaxWidth()
                        .defaultMinSize(minHeight = 56.dp),
                    label = {
                        FjarmNormalSizeText(
                            text = passwordInputLabelText,
                            fontWeight = FontWeight.Bold,
                        )
                    },
                    visualTransformation = PasswordVisualTransformation(),
                    shape = RoundedCornerShape(4.dp),
                    keyboardOptions = KeyboardOptions(
                        keyboardType = KeyboardType.Password,
                        imeAction = ImeAction.Done,
                    ),
                    singleLine = true,
                )

                Spacer(modifier = Modifier.height(24.dp))

                // Primary Action Button
                FjarmFilledButton(
                    onClick = { },
                    modifier = Modifier
                        .fillMaxWidth()
                        .defaultMinSize(minHeight = 56.dp),
                    enabled = loginButtonEnabled,
                ) {
                    FjarmNormalSizeText(
                        text = loginButtonText,
                        fontWeight = FontWeight.Black,
                    )
                }

                Spacer(modifier = Modifier.height(48.dp))

                // Divider
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    HorizontalDivider(
                        modifier = Modifier.weight(1f),
                    )
                    FjarmNormalSizeText(
                        text = alternativeOptionsText,
                        modifier = Modifier.padding(horizontal = 16.dp),
                        fontWeight = FontWeight.Bold
                    )
                    HorizontalDivider(
                        modifier = Modifier.weight(1f),
                    )
                }

                Spacer(modifier = Modifier.height(48.dp))

                // Footer Links
                FjarmNormalSizeText(
                    text = newToFjarmPromptText,
                    fontWeight = FontWeight.Bold,
                )
                FjarmButton(
                    modifier = Modifier
                        .wrapContentHeight(),
                    onClick = { },
                ) {
                    FjarmNormalSizeText(
                        text = navigateToSignUpButtonText,
                        fontWeight = FontWeight.Bold,
                        textDecoration = TextDecoration.Underline,
                    )
                }

                Spacer(modifier = Modifier.height(40.dp))

                // Legal Links
                Row(
                    modifier = Modifier
                        .fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceEvenly,
                ) {
                    FjarmSmallSizeText(privacyPolicyText)
                    FjarmSmallSizeText(termsOfServiceText)
                }
            }
        }
    }
}

@PreviewFontScale
@PreviewLightDark
@PreviewScreenSizes
@Composable
fun LoginScreenContentPreview() {
    LoginScreenContent()
}
