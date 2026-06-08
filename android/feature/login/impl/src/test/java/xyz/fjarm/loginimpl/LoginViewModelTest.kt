package xyz.fjarm.loginimpl

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
import org.junit.Before
import org.junit.Test
import xyz.fjarm.loginlibrary.AttemptLoginUseCase

@OptIn(ExperimentalCoroutinesApi::class)
class LoginViewModelTest {

    private val noopAttemptLoginUseCase = object: AttemptLoginUseCase {
        override suspend fun invoke(email: String, password: String): Result<Unit> {
            TODO("Not yet implemented")
        }
    }

    private val successAttemptLoginUseCase = object: AttemptLoginUseCase {
        override suspend fun invoke(email: String, password: String): Result<Unit> {
            return Result.success(Unit)
        }
    }

    private val failureAttemptLoginUseCase = object: AttemptLoginUseCase {
        override suspend fun invoke(email: String, password: String): Result<Unit> {
            return Result.failure(Exception("Something went wrong."))
        }

    }

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
    fun processEvent_EmailAddressModified_emitsNoSideEffects_andUpdatesEmailState() = runTest {
        // Given a LoginViewModel with a no-op AttemptLoginUseCase
        val viewModel = LoginViewModel(noopAttemptLoginUseCase)

        val collectedSideEffects = mutableListOf<LoginSideEffect>()
        backgroundScope.launch(UnconfinedTestDispatcher(testScheduler)) {
            viewModel.sideEffect.collect { collectedSideEffects.add(it) }
        }

        // When a user modifies the email edit text
        viewModel.processEvent(LoginEvent.EmailAddressModified("j"))
        advanceUntilIdle()

        // Then no side effects are emitted
        assertEquals(0, collectedSideEffects.size)

        // And the email state is updated
        assertEquals("j", viewModel.state.value.userInput.emailInputText)
        assertEquals(true, viewModel.state.value.userInput.emailInputIsInvalid)

        // And the login button is not enabled
        assertEquals(false, viewModel.state.value.loginButton.loginButtonEnabled)
    }

    @Test
    fun processEvent_EmailAddressModified_multipleUpdates_updatesEmailState() = runTest {
        // Given a LoginViewModel with a no-op AttemptLoginUseCase
        val viewModel = LoginViewModel(noopAttemptLoginUseCase)

        // When a user modifies the email edit text
        viewModel.processEvent(LoginEvent.EmailAddressModified("j"))
        advanceUntilIdle()
        viewModel.processEvent(LoginEvent.EmailAddressModified("ja"))
        advanceUntilIdle()

        // The email state is updated
        assertEquals("ja", viewModel.state.value.userInput.emailInputText)
        assertEquals(true, viewModel.state.value.userInput.emailInputIsInvalid)

        // And the login button is not enabled
        assertEquals(false, viewModel.state.value.loginButton.loginButtonEnabled)

        viewModel.processEvent(LoginEvent.EmailAddressModified("ja@d.co"))
        advanceUntilIdle()

        // Then the email state is updated
        assertEquals("ja@d.co", viewModel.state.value.userInput.emailInputText)
        assertEquals(false, viewModel.state.value.userInput.emailInputIsInvalid)

        // And the login button is not enabled
        assertEquals(false, viewModel.state.value.loginButton.loginButtonEnabled)
    }

    @Test
    fun processEvent_PasswordModified_emitsNoSideEffects_andUpdatesPasswordState() = runTest {
        val viewModel = LoginViewModel(noopAttemptLoginUseCase)

        val collectedSideEffects = mutableListOf<LoginSideEffect>()
        backgroundScope.launch(UnconfinedTestDispatcher(testScheduler)) {
            viewModel.sideEffect.collect { collectedSideEffects.add(it) }
        }

        // When a user modifies the password edit text
        viewModel.processEvent(LoginEvent.PasswordModified("j"))
        advanceUntilIdle()

        // Then no side effects are emitted
        assertEquals(0, collectedSideEffects.size)

        // And the password state is updated
        assertEquals("j", viewModel.state.value.userInput.passwordInputText)

        // And the login button is not enabled
        assertEquals(false, viewModel.state.value.loginButton.loginButtonEnabled)
    }

    @Test
    fun processEvent_EmailAddressModified_PasswordModified_validNonEmptyValuesUpdatesLoginButtonState() = runTest {
        // Given a LoginViewModel with a no-op AttemptLoginUseCase
        val viewModel = LoginViewModel(noopAttemptLoginUseCase)

        // When a user modifies the email and password edit texts with valid non-empty values
        viewModel.processEvent(LoginEvent.EmailAddressModified("j@d.co"))
        advanceUntilIdle()
        viewModel.processEvent(LoginEvent.PasswordModified("bleepbloop"))
        advanceUntilIdle()

        // Then the login button is enabled
        assertEquals(true, viewModel.state.value.loginButton.loginButtonEnabled)
    }
}