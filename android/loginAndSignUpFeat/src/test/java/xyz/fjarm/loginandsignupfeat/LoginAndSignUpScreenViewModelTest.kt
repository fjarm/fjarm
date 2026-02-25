package xyz.fjarm.loginandsignupfeat

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.launch
import kotlinx.coroutines.test.StandardTestDispatcher
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.advanceUntilIdle
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test

@OptIn(ExperimentalCoroutinesApi::class)
class LoginAndSignUpScreenViewModelTest {

    private val testDispatcher = StandardTestDispatcher()

    @Before
    fun setUp() {
        Dispatchers.setMain(testDispatcher)
    }

    @After
    fun tearDown() {
        Dispatchers.resetMain()
    }

    @Test
    fun processEvent_SignUpButtonClicked_emitsNavigateToSignUpSideEffect() = runTest {
        // Given a LoginAndSignUpScreenViewModel
        val viewModel = LoginAndSignUpScreenViewModel()

        val collectedSideEffects = mutableListOf<LoginAndSignUpSideEffect>()
        backgroundScope.launch(UnconfinedTestDispatcher(testScheduler)) {
            viewModel.sideEffect.collect { collectedSideEffects.add(it) }
        }

        // When the SignUpButtonClicked event is processed
        viewModel.processEvent(LoginAndSignUpEvent.SignUpButtonClicked)
        advanceUntilIdle()

        // Then a NavigateToSignUp side effect is emitted
        val sideEffect = collectedSideEffects.firstOrNull()
        assertEquals(1, collectedSideEffects.size)
        assertTrue(sideEffect is LoginAndSignUpSideEffect.NavigateToSignUp)
    }

    @Test
    fun processEvent_LoginButtonClicked_emitsNavigateToLoginSideEffect() = runTest {
        // Given a LoginAndSignUpScreenViewModel
        val viewModel = LoginAndSignUpScreenViewModel()

        val collectedSideEffects = mutableListOf<LoginAndSignUpSideEffect>()
        backgroundScope.launch(UnconfinedTestDispatcher(testScheduler)) {
            viewModel.sideEffect.collect { collectedSideEffects.add(it) }
        }

        // When the SignUpButtonClicked event is processed
        viewModel.processEvent(LoginAndSignUpEvent.LogInButtonClicked)
        advanceUntilIdle()

        // Then a NavigateToSignUp side effect is emitted
        val sideEffect = collectedSideEffects.firstOrNull()
        assertEquals(1, collectedSideEffects.size)
        assertTrue(sideEffect is LoginAndSignUpSideEffect.NavigateToLogIn)
    }
}
