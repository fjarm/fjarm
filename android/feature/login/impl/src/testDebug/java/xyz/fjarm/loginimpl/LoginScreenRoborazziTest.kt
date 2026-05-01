package xyz.fjarm.loginimpl

import androidx.compose.ui.test.junit4.v2.createAndroidComposeRule
import androidx.compose.ui.test.onRoot
import com.github.takahirom.roborazzi.RobolectricDeviceQualifiers
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

@OptIn(ExperimentalCoroutinesApi::class)
@HiltAndroidTest
@Config(
    application = HiltTestApplication::class,
    sdk = [36],
    qualifiers = RobolectricDeviceQualifiers.Pixel6,
)
@GraphicsMode(GraphicsMode.Mode.NATIVE)
@RunWith(RobolectricTestRunner::class)
class LoginScreenRoborazziTest {

    @get:Rule(order = 0)
    var hiltRule = HiltAndroidRule(this)

    @get:Rule(order = 1)
    val composeTestRule = createAndroidComposeRule<TestActivity>()

    @Before
    fun setUp() {
        hiltRule.inject()
    }

    @Test
    @Config(qualifiers = "+night")
    fun loginScreenTest_darkMode() = runTest {
        composeTestRule.setContent {
            LoginScreen()
        }

        composeTestRule.waitForIdle()
        composeTestRule
            .onRoot()
            .captureRoboImage()
    }

    @Test
    @Config(qualifiers = "+notnight")
    fun loginScreenTest_lightMode() = runTest {
        composeTestRule.setContent {
            LoginScreen()
        }

        composeTestRule.waitForIdle()
        composeTestRule
            .onRoot()
            .captureRoboImage()
    }
}
