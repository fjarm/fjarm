package xyz.fjarm.loginimpl

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
import androidx.compose.foundation.layout.systemBarsPadding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.layout.wrapContentHeight
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
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
import xyz.fjarm.buttons.FjarmButton
import xyz.fjarm.buttons.FjarmFilledButton
import xyz.fjarm.fjarmtheme.FjarmTheme

@Composable
fun LoginScreenContent(
    modifier: Modifier = Modifier
) {
    FjarmTheme {
        Surface(
            modifier = modifier
                .fillMaxSize(),
        ) {
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .systemBarsPadding()
                    .padding(horizontal = 32.dp, vertical = 16.dp),
                horizontalAlignment = Alignment.CenterHorizontally,
            ) {
                // Brand Logo
                Text(
                    text = "Fjarm",
                    fontSize = 24.sp,
                    fontWeight = FontWeight.Black,
                    modifier = Modifier.align(Alignment.Start),
                )

                Spacer(modifier = Modifier.height(80.dp))

                // Large Header
                Column(modifier = Modifier.fillMaxWidth()) {
                    Text(
                        text = "Login",
                        fontSize = 80.sp,
                        fontWeight = FontWeight.Black,
                        lineHeight = 72.sp,
                        letterSpacing = (-2).sp,
                    )
                    Box(
                        modifier = Modifier
                            .width(124.dp)
                            .height(4.dp)
                            .background(Color.Black),
                    )
                }

                Spacer(modifier = Modifier.height(64.dp))

                // Input Field Section
                Column(modifier = Modifier.fillMaxWidth()) {
                    Text(
                        text = "Email address",
                        fontSize = 12.sp,
                        fontWeight = FontWeight.Bold,
                        letterSpacing = 1.sp,
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                    OutlinedTextField(
                        value = "",
                        onValueChange = { },
                        modifier = Modifier
                            .fillMaxWidth()
                            .height(56.dp),
                        placeholder = {
                            Text(
                                text = "name@domain.com",
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
                }

                Spacer(modifier = Modifier.height(32.dp))

                // Primary Action Button
                FjarmFilledButton(
                    onClick = { },
                    modifier = Modifier
                        .fillMaxWidth()
                        .height(56.dp),
                ) {
                    Text(
                        text = "Continue",
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
                    HorizontalDivider(
                        modifier = Modifier.weight(1f),
                    )
                    Text(
                        text = "Or",
                        modifier = Modifier.padding(horizontal = 16.dp),
                        fontSize = 12.sp,
                        fontWeight = FontWeight.Bold
                    )
                    HorizontalDivider(
                        modifier = Modifier.weight(1f),
                    )
                }

                Spacer(modifier = Modifier.height(48.dp))

                // Footer Links
                Column(horizontalAlignment = Alignment.CenterHorizontally) {
                    Text(
                        text = "New to Fjarm?",
                        fontSize = 12.sp,
                        fontWeight = FontWeight.Bold,
                        letterSpacing = 1.sp,
                    )
                    FjarmButton(
                        modifier = Modifier
                            .wrapContentHeight(),
                        onClick = { },
                    ) {
                        Text(
                            text = "Sign up",
                            fontSize = 16.sp,
                            fontWeight = FontWeight.Black,
                            textDecoration = TextDecoration.Underline,
                        )
                    }
                }

                Spacer(modifier = Modifier.height(40.dp))

                // Legal Links
                Row(
                    modifier = Modifier
                        .fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceEvenly,
                ) {
                    LegalLink("Privacy policy")
                    LegalLink("Terms of service")
                }
            }
        }
    }
}

@Composable
fun LegalLink(text: String) {
    Text(
        text = text,
        fontSize = 10.sp,
        fontWeight = FontWeight.Bold,
        letterSpacing = 0.5.sp,
        textAlign = TextAlign.Center,
    )
}

@PreviewFontScale
@PreviewLightDark
@PreviewScreenSizes
@Composable
fun LoginScreenContentPreview() {
    LoginScreenContent()
}
