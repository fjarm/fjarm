package xyz.fjarm.loginimpl

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.OutlinedTextFieldDefaults
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.tooling.preview.PreviewFontScale
import androidx.compose.ui.tooling.preview.PreviewLightDark
import androidx.compose.ui.tooling.preview.PreviewScreenSizes
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

/**
 * MVI State for the Login Email Screen
 */
data class LoginEmailUiState(
    val email: String = "",
    val isLoading: Boolean = false,
    val error: String? = null
)

/**
 * MVI Intent for the Login Email Screen
 */
sealed class LoginEmailIntent {
    data class EmailChanged(val email: String) : LoginEmailIntent()
    object ContinueClicked : LoginEmailIntent()
    object GoogleLoginClicked : LoginEmailIntent()
    object AppleLoginClicked : LoginEmailIntent()
    object CreateAccountClicked : LoginEmailIntent()
}

//@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun LoginEmailScreen(
    state: LoginEmailUiState,
    onIntent: (LoginEmailIntent) -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .fillMaxSize()
            .background(Color(0xFFFAFAF5))
            .padding(horizontal = 24.dp, vertical = 40.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        // Brand Logo
        Text(
            text = "MONOLITH",
            fontSize = 24.sp,
            fontWeight = FontWeight.Black,
            letterSpacing = (-1).sp,
            modifier = Modifier.align(Alignment.Start)
        )

        Spacer(modifier = Modifier.height(80.dp))

        // Large Header
        Column(modifier = Modifier.fillMaxWidth()) {
            Text(
                text = "LOGIN",
                fontSize = 80.sp,
                fontWeight = FontWeight.Black,
                lineHeight = 72.sp,
                letterSpacing = (-4).sp
            )
            Box(
                modifier = Modifier
                    .width(100.dp)
                    .height(4.dp)
                    .background(Color.Black)
                    .padding(top = 8.dp)
            )
        }

        Spacer(modifier = Modifier.height(64.dp))

        // Input Field Section
        Column(modifier = Modifier.fillMaxWidth()) {
            Text(
                text = "EMAIL ADDRESS",
                fontSize = 12.sp,
                fontWeight = FontWeight.Bold,
                letterSpacing = 1.sp,
                color = Color.DarkGray
            )
            Spacer(modifier = Modifier.height(8.dp))
            OutlinedTextField(
                value = state.email,
                onValueChange = { onIntent(LoginEmailIntent.EmailChanged(it)) },
                modifier = Modifier
                    .fillMaxWidth()
                    .height(64.dp),
                placeholder = {
                    Text(
                        text = "NAME@DOMAIN.COM",
                        color = Color.LightGray,
                        fontWeight = FontWeight.Bold
                    )
                },
                colors = OutlinedTextFieldDefaults.colors(
                    focusedBorderColor = Color.Black,
                    unfocusedBorderColor = Color.Black,
                    cursorColor = Color.Black
                ),
                shape = RoundedCornerShape(0.dp),
                keyboardOptions = KeyboardOptions(
                    keyboardType = KeyboardType.Email,
                    imeAction = ImeAction.Next
                ),
                singleLine = true
            )
        }

        Spacer(modifier = Modifier.height(32.dp))

        // Primary Action Button
        Button(
            onClick = { onIntent(LoginEmailIntent.ContinueClicked) },
            modifier = Modifier
                .fillMaxWidth()
                .height(72.dp),
            colors = ButtonDefaults.buttonColors(containerColor = Color.Black),
            shape = RoundedCornerShape(0.dp)
        ) {
            Text(
                text = "CONTINUE",
                color = Color.White,
                fontSize = 18.sp,
                fontWeight = FontWeight.Black,
                letterSpacing = 2.sp
            )
        }

        Spacer(modifier = Modifier.height(48.dp))

        // Divider
        Row(
            modifier = Modifier.fillMaxWidth(),
            verticalAlignment = Alignment.CenterVertically
        ) {
            HorizontalDivider(modifier = Modifier.weight(1f), color = Color.LightGray)
            Text(
                text = "OR",
                modifier = Modifier.padding(horizontal = 16.dp),
                fontSize = 10.sp,
                color = Color.Gray,
                fontWeight = FontWeight.Bold
            )
            HorizontalDivider(modifier = Modifier.weight(1f), color = Color.LightGray)
        }

        Spacer(modifier = Modifier.height(48.dp))

        // Social Login Buttons
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            SocialButton(
                text = "GOOGLE",
                onClick = { onIntent(LoginEmailIntent.GoogleLoginClicked) },
                modifier = Modifier.weight(1f)
            )
            SocialButton(
                text = "APPLE",
                onClick = { onIntent(LoginEmailIntent.AppleLoginClicked) },
                modifier = Modifier.weight(1f)
            )
        }

        Spacer(modifier = Modifier.weight(1f))

        // Footer Links
        Column(horizontalAlignment = Alignment.CenterHorizontally) {
            Text(
                text = "NEW TO THE ARCHITECTURAL MONOLITH?",
                fontSize = 10.sp,
                fontWeight = FontWeight.Bold,
                letterSpacing = 1.sp,
                color = Color.Gray
            )
            TextButton(onClick = { onIntent(LoginEmailIntent.CreateAccountClicked) }) {
                Text(
                    text = "CREATE AN ACCOUNT",
                    fontSize = 12.sp,
                    fontWeight = FontWeight.Black,
                    color = Color.Black,
                    textDecoration = TextDecoration.Underline
                )
            }
        }

        Spacer(modifier = Modifier.height(40.dp))

        // Legal Links
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween
        ) {
            LegalLink("PRIVACY POLICY")
            LegalLink("TERMS OF SERVICE")
            LegalLink("LEGAL NOTICE")
        }
    }
}

@Composable
fun SocialButton(
    text: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    OutlinedButton(
        onClick = onClick,
        modifier = modifier.height(56.dp),
        shape = RoundedCornerShape(0.dp),
        border = BorderStroke(1.dp, Color.Black),
        colors = ButtonDefaults.outlinedButtonColors(contentColor = Color.Black)
    ) {
        Text(
            text = text,
            fontSize = 12.sp,
            fontWeight = FontWeight.Black,
            letterSpacing = 1.sp
        )
    }
}

@Composable
fun LegalLink(text: String) {
    Text(
        text = text,
        fontSize = 10.sp,
        fontWeight = FontWeight.Bold,
        letterSpacing = 0.5.sp,
        color = Color.Gray,
        textAlign = TextAlign.Center
    )
}

@PreviewFontScale
@PreviewLightDark
@PreviewScreenSizes
@Composable
fun LoginEmailScreenPreview() {
    LoginEmailScreen(
        state = LoginEmailUiState(email = ""),
        onIntent = {}
    )
}
