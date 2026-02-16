package xyz.fjarm.loginandsignupfeat

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp
import xyz.fjarm.fjarmthemelib.FjarmTheme
import xyz.fjarm.previewslib.PreviewFontScales
import xyz.fjarm.previewslib.PreviewLightDarkTheme
import xyz.fjarm.previewslib.PreviewPhoneSizes

@Composable
fun LoginAndSignUpScreen(
    // TODO: Add ViewModel dependency to replace mutableStateOf calls below
    modifier: Modifier = Modifier,
) {
    FjarmTheme {

        Column(modifier = modifier.padding(16.dp)) {
            // Username Field
            val username = remember { mutableStateOf("") }
            val password = remember { mutableStateOf("") }
            val showUsernameError = remember { mutableStateOf(false) }
            val showPasswordError = remember { mutableStateOf(false) }

            // Username Field
            OutlinedTextField(
                value = username.value,
                onValueChange = {
                    username.value = it
                    showUsernameError.value = it.isBlank()
                },
                label = { Text("Username") },
                isError = showUsernameError.value,
                shape = RoundedCornerShape(8.dp),
                modifier = Modifier
                    .fillMaxWidth(),
            )
            if (showUsernameError.value) {
                Text(text = "Username cannot be empty", color = MaterialTheme.colorScheme.error)
            }

            Spacer(modifier = Modifier.height(16.dp))

            // Password Field
            OutlinedTextField(
                value = password.value,
                onValueChange = {
                    password.value = it
                    showPasswordError.value = it.length < 8
                },
                label = { Text("Password") },
                isError = showPasswordError.value,
                visualTransformation = PasswordVisualTransformation(),
                shape = RoundedCornerShape(8.dp),
                modifier = Modifier
                    .fillMaxWidth(),
            )
            if (showPasswordError.value) {
                Text("Password must be at least 8 characters", color = MaterialTheme.colorScheme.error)
            }

            Spacer(modifier = Modifier.height(16.dp))

            // Submit Button
            Button(
                onClick = { /* TODO: Handle login */ },
                enabled = username.value.isNotBlank() && password.value.isNotBlank() && !showUsernameError.value && !showPasswordError.value,
                shape = RoundedCornerShape(8.dp),
                modifier = Modifier
                    .fillMaxWidth(),
            ) {
                Text("Log In")
            }
        }
    }
}

@PreviewPhoneSizes
@PreviewFontScales
@PreviewLightDarkTheme
@Composable
fun LoginAndSignUpScreenPreview() {
    LoginAndSignUpScreen(
        modifier = Modifier,
    )
}
