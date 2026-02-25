package xyz.fjarm.loginandsignupfeat

import androidx.compose.ui.test.junit4.createAndroidComposeRule
import androidx.compose.ui.test.onNodeWithText
import androidx.compose.ui.test.performClick
import androidx.hilt.lifecycle.viewmodel.compose.hiltViewModel
import dagger.hilt.android.testing.HiltAndroidRule
import dagger.hilt.android.testing.HiltAndroidTest
import dagger.hilt.android.testing.HiltTestApplication
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.launch
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.advanceUntilIdle
import kotlinx.coroutines.test.runTest
import org.junit.Assert
import org.junit.Before
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith
import org.robolectric.RobolectricTestRunner
import org.robolectric.annotation.Config
import xyz.fjarm.testactivitylib.TestActivity

@OptIn(ExperimentalCoroutinesApi::class)
@HiltAndroidTest
@Config(application = HiltTestApplication::class, sdk = [31, 36])
@RunWith(RobolectricTestRunner::class)
class LoginAndSignUpScreenTest {

    @get:Rule(order = 0)
    var hiltRule = HiltAndroidRule(this)

    @get:Rule(order = 1)
    val composeTestRule = createAndroidComposeRule<TestActivity>()

    @Before
    fun setUp() {
        hiltRule.inject()
    }

    @Test
    fun logInButtonClick_emitsNavigateToLogInSideEffect() = runTest {
        // Given the Login and Sign Up screen is displayed
        val logInText = composeTestRule.activity.getString(R.string.log_in_button)

        val collectedSideEffects = mutableListOf<LoginAndSignUpSideEffect>()
        lateinit var viewModel: LoginAndSignUpScreenViewModel

        composeTestRule.setContent {
            viewModel = hiltViewModel<LoginAndSignUpScreenViewModel>()
            LoginAndSignUpScreen(
                viewModel = viewModel
            )
        }

        // Start collecting side effects in the background
        backgroundScope.launch(UnconfinedTestDispatcher(testScheduler)) {
            viewModel.sideEffect.collect { collectedSideEffects.add(it) }
        }

        // When the log in button is clicked
        composeTestRule.onNodeWithText(logInText).performClick()
        advanceUntilIdle()

        // Then a NavigateToLogIn side effect is received
        val sideEffect = collectedSideEffects.firstOrNull()
        Assert.assertEquals(1, collectedSideEffects.size)
        Assert.assertTrue(
            "Expected NavigateToSignUp side effect but received $sideEffect",
            sideEffect is LoginAndSignUpSideEffect.NavigateToLogIn
        )
    }

    @Test
    fun signUpButtonClick_emitsNavigateToSignUpSideEffect() = runTest {
        // Given the Login and Sign Up screen is displayed
        val signUpText = composeTestRule.activity.getString(R.string.sign_up_button)

        val collectedSideEffects = mutableListOf<LoginAndSignUpSideEffect>()
        lateinit var viewModel: LoginAndSignUpScreenViewModel

        composeTestRule.setContent {
            viewModel = hiltViewModel<LoginAndSignUpScreenViewModel>()
            LoginAndSignUpScreen(
                viewModel = viewModel
            )
        }

        // Start collecting side effects in the background
        backgroundScope.launch(UnconfinedTestDispatcher(testScheduler)) {
            viewModel.sideEffect.collect { collectedSideEffects.add(it) }
        }

        // When the sign up button is clicked
        composeTestRule.onNodeWithText(signUpText).performClick()
        advanceUntilIdle()

        // Then a NavigateToSignUp side effect is received
        val sideEffect = collectedSideEffects.firstOrNull()
        Assert.assertEquals(1, collectedSideEffects.size)
        Assert.assertTrue(
            "Expected NavigateToSignUp side effect but received $sideEffect",
            sideEffect is LoginAndSignUpSideEffect.NavigateToSignUp
        )
    }
}