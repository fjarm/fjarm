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
import androidx.compose.foundation.layout.systemBarsPadding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.layout.wrapContentHeight
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Surface
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.tooling.preview.PreviewFontScale
import androidx.compose.ui.tooling.preview.PreviewLightDark
import androidx.compose.ui.tooling.preview.PreviewScreenSizes
import androidx.compose.ui.unit.dp
import xyz.fjarm.buttons.FjarmButton
import xyz.fjarm.buttons.FjarmFilledButton
import xyz.fjarm.fjarmtheme.FjarmTheme
import xyz.fjarm.text.FjarmExtraLargeHeaderText
import xyz.fjarm.text.FjarmNormalSizeText
import xyz.fjarm.text.FjarmSmallSizeText

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

                Spacer(modifier = Modifier.height(80.dp))

                // Large Header
                Column(
                    modifier = Modifier
                        .fillMaxWidth(),
                ) {
                    FjarmExtraLargeHeaderText(
                        text = "Login",
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
                    value = "",
                    onValueChange = { },
                    modifier = Modifier
                        .fillMaxWidth()
                        .defaultMinSize(minHeight = 56.dp),
                    label = {
                        FjarmNormalSizeText(
                            text = "Email address",
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
                    value = "",
                    onValueChange = { },
                    modifier = Modifier
                        .fillMaxWidth()
                        .defaultMinSize(minHeight = 56.dp),
                    label = {
                        FjarmNormalSizeText(
                            text = "Password",
                            fontWeight = FontWeight.Bold,
                        )
                    },
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
                ) {
                    FjarmNormalSizeText(
                        text = "Continue",
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
                        text = "Or",
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
                    text = "New to Fjarm?",
                    fontWeight = FontWeight.Bold,
                )
                FjarmButton(
                    modifier = Modifier
                        .wrapContentHeight(),
                    onClick = { },
                ) {
                    FjarmNormalSizeText(
                        text = "Sign up",
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
                    FjarmSmallSizeText("Privacy policy")
                    FjarmSmallSizeText("Terms of service")
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
