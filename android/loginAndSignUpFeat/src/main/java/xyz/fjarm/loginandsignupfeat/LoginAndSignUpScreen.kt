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
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import xyz.fjarm.fjarmthemelib.FjarmTheme

@Composable
fun LoginAndSignUpScreen(
    modifier: Modifier = Modifier,
) {
    FjarmTheme {

        Column(modifier = Modifier.padding(16.dp)) {
            // Username Field
            OutlinedTextField(
                value = "Bleep bloop",
                onValueChange = {  },
                label = { Text("Username") },
                isError = false, // Highlights the text field in red if there's an error. Should align with the error below
                shape = RoundedCornerShape(8.dp),
                modifier = Modifier
                    .fillMaxWidth(),
            )
            if (true) { // If there's username error
                Text(text = "Oh no there's an error", color = MaterialTheme.colorScheme.error)
            }

            Spacer(modifier = Modifier.height(16.dp))

            // Password Field
            OutlinedTextField(
                value = "Gleep gloop",
                onValueChange = {  },
                label = { Text("Password") },
                isError = false,
                visualTransformation = PasswordVisualTransformation(),
                shape = RoundedCornerShape(8.dp),
                modifier = Modifier
                    .fillMaxWidth(),
            )
            if (true) {
                Text("Oh no there's an error", color = MaterialTheme.colorScheme.error)
            }

            Spacer(modifier = Modifier.height(16.dp))

            // Submit Button
            Button(
                onClick = {  },
                enabled = true,
                shape = RoundedCornerShape(8.dp),
                modifier = Modifier
                    .fillMaxWidth(),
            ) {
                Text("Log In")
            }
        }
    }
}

@Preview(showBackground = true)
@Composable
fun LoginAndSignUpScreenPreview() {
    LoginAndSignUpScreen(
        modifier = Modifier,
    )
}
