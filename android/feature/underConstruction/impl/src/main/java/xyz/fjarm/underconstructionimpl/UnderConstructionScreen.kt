package xyz.fjarm.underconstructionimpl

import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.systemBarsPadding
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import xyz.fjarm.fjarmtheme.FjarmTheme
import xyz.fjarm.fjarmtheme.FjarmTypography
import xyz.fjarm.previews.PreviewFontScales
import xyz.fjarm.previews.PreviewLightDarkTheme
import xyz.fjarm.previews.PreviewPhoneSizes

@Composable
fun UnderConstructionScreen() {
    FjarmTheme {
        Surface(
            modifier = Modifier
                .fillMaxSize(),
        ) {
            Column(
                modifier = Modifier
                    .fillMaxSize()
                    .systemBarsPadding()
                    .padding(horizontal = 32.dp, vertical = 16.dp),
                horizontalAlignment = Alignment.CenterHorizontally,
                // SpacedBy ensures the header and image have distance if the screen is small
                verticalArrangement = Arrangement.spacedBy(16.dp),
            ) {
                // 1. Header
                Text(
                    text = "This feature is under construction",
                    modifier = Modifier
                        .fillMaxWidth(),
                    style = FjarmTypography.headlineLarge.copy(
                        fontWeight = FontWeight.Bold,
                    ),
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
                        painter = painterResource(id = R.drawable.construction_worker),
                        contentDescription = null,
                        modifier = Modifier.fillMaxSize(),
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
fun UnderConstructionScreenPreview() {
    UnderConstructionScreen()
}
