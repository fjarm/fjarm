package xyz.fjarm.buttons

import androidx.compose.material3.Text
import androidx.compose.ui.test.junit4.createAndroidComposeRule
import androidx.compose.ui.test.onRoot
import com.github.takahirom.roborazzi.captureRoboImage
import dagger.hilt.android.testing.HiltAndroidRule
import dagger.hilt.android.testing.HiltAndroidTest
import dagger.hilt.android.testing.HiltTestApplication
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.runTest
import org.junit.Before
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith
import org.robolectric.RobolectricTestRunner
import org.robolectric.annotation.Config
import org.robolectric.annotation.GraphicsMode
import xyz.fjarm.testactivity.TestActivity
import kotlin.intArrayOf

@OptIn(ExperimentalCoroutinesApi::class)
@HiltAndroidTest
@Config(application = HiltTestApplication::class, sdk = [36])
@GraphicsMode(GraphicsMode.Mode.NATIVE)
@RunWith(RobolectricTestRunner::class)
class FjarmButtonsRoborazziTest {

    @get:Rule(order = 0)
    var hiltRule = HiltAndroidRule(this)

    @get:Rule(order = 1)
    val composeTestRule = createAndroidComposeRule<TestActivity>()

    @Before
    fun setUp() {
        hiltRule.inject()
    }

    @Test
    fun button_lightMode() = runTest {
        composeTestRule.setContent {
            FjarmButton(onClick = {}) {
                Text("Button")
            }
        }

        composeTestRule.waitForIdle()
        composeTestRule
            .onRoot()
            .captureRoboImage()
    }

    @Test
    fun outlinedButton_lightMode() = runTest {
        composeTestRule.setContent {
            FjarmOutlinedButton(onClick = {}) {
                Text("Button")
            }
        }

        composeTestRule.waitForIdle()
        composeTestRule
            .onRoot()
            .captureRoboImage()
    }
}
