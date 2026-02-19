package xyz.fjarm.loginandsignupfeat

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.padding
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.hilt.lifecycle.viewmodel.compose.hiltViewModel
import xyz.fjarm.fjarmthemelib.FjarmTheme
import xyz.fjarm.previewslib.PreviewFontScales
import xyz.fjarm.previewslib.PreviewLightDarkTheme
import xyz.fjarm.previewslib.PreviewPhoneSizes

@Composable
fun LoginAndSignUpScreen(
    modifier: Modifier = Modifier,
    viewModel: LoginAndSignUpScreenViewModel = hiltViewModel(),
) {
    FjarmTheme {
        Column(modifier = modifier.padding(16.dp)) {
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
