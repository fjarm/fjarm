package xyz.fjarm.navigationlib

import androidx.lifecycle.SavedStateHandle
import androidx.navigation3.runtime.NavKey
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.StandardTestDispatcher
import kotlinx.coroutines.test.advanceUntilIdle
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Before
import org.junit.Test
import xyz.fjarm.loginandsignupfeatlib.LoginAndSignUpNavKey

@OptIn(ExperimentalCoroutinesApi::class)
class NavigatorImplTest {

    private val testDispatcher = StandardTestDispatcher()

    @Before
    fun setUp() {
        Dispatchers.setMain(testDispatcher)
    }

    @After
    fun tearDown() {
        Dispatchers.resetMain()
    }

    // No [Serializable] annotation is needed because SavedStateHandle is backed by an in-memory map
    // in unit tests.
    private data object DummyDestinationNavKey : NavKey

    @Test
    fun navigateTo_withDummyDestination_addsDummyDestinationToBackStack() = runTest {
        // Given a fresh NavigatorImpl whose back stack contains only LoginAndSignUpNavKey
        val navigator = NavigatorImpl(SavedStateHandle())
        advanceUntilIdle()

        // When navigateTo is called with DummyDestinationNavKey
        navigator.navigateTo(DummyDestinationNavKey)
        advanceUntilIdle()

        // Then the back stack contains LoginAndSignUpNavKey followed by DummyDestinationNavKey
        val backStack = navigator.getBackStack()
        assertEquals(2, backStack.size)
        assertEquals(LoginAndSignUpNavKey, backStack[0])
        assertEquals(DummyDestinationNavKey, backStack[1])
    }
}
