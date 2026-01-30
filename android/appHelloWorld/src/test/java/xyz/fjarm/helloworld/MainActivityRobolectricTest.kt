package xyz.fjarm.helloworld

import android.widget.Toast
import androidx.compose.ui.test.junit4.createAndroidComposeRule
import androidx.compose.ui.test.onNodeWithText
import androidx.compose.ui.test.performClick
import build.buf.gen.fjarm.helloworld.v1.HelloWorldOutput
import dagger.hilt.android.testing.BindValue
import dagger.hilt.android.testing.HiltAndroidRule
import dagger.hilt.android.testing.HiltAndroidTest
import dagger.hilt.android.testing.HiltTestApplication
import dagger.hilt.android.testing.UninstallModules
import org.junit.Assert.assertEquals
import org.junit.Before
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith
import org.robolectric.RobolectricTestRunner
import org.robolectric.annotation.Config
import org.robolectric.shadows.ShadowToast
import xyz.fjarm.libhelloworld.GetHelloWorldUseCase
import xyz.fjarm.libhelloworld.GetHelloWorldUseCase.GetHelloWorldUseCaseModule

@HiltAndroidTest
@UninstallModules(GetHelloWorldUseCaseModule::class)
@Config(application = HiltTestApplication::class, sdk = [31, 36])
@RunWith(RobolectricTestRunner::class)
class MainActivityRobolectricTest {

    @get:Rule(order = 0)
    var hiltRule = HiltAndroidRule(this)

    @get:Rule(order = 1)
    val composeTestRule = createAndroidComposeRule<MainActivity>()

    private val expectedMessage = "Hello from Fake Use Case!"

    @BindValue
    @JvmField
    val fakeGetHelloWorldUseCase: GetHelloWorldUseCase = object : GetHelloWorldUseCase {
        override suspend fun invoke(): HelloWorldOutput {
            return HelloWorldOutput.newBuilder()
                .setOutput(expectedMessage)
                .build()
        }
    }

    @Before
    fun init() {
        hiltRule.inject()
    }

    @Test
    fun buttonClick_withoutViewModelFailure_showsToast() {
        // Given a visible button that can be clicked
        val buttonText = composeTestRule.activity.getString(R.string.button_text)

        // When the button is clicked
        composeTestRule.onNodeWithText(buttonText).performClick()

        // We wait for the side effect to be processed.
        // Robolectric's main looper will execute the LaunchedEffect and emission.
        composeTestRule.waitForIdle()

        val latestToast = ShadowToast.getLatestToast()
        val toastMessage = ShadowToast.getTextOfLatestToast()

        // Then a toast is shown with the expected message
        assert(latestToast != null)
        assertEquals(expectedMessage, toastMessage)
        assertEquals(Toast.LENGTH_SHORT, latestToast.duration)
    }
}
