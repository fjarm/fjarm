package xyz.fjarm.navigationlib

import android.os.Parcelable
import androidx.lifecycle.SavedStateHandle
import androidx.navigation3.runtime.NavKey
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.StandardTestDispatcher
import kotlinx.coroutines.test.advanceUntilIdle
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import kotlinx.parcelize.Parcelize
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertThrows
import org.junit.Before
import org.junit.Test

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

    @Parcelize
    private data object ValidStartDestinationNavKey: NavKey, Parcelable

    // No [Serializable] annotation is needed because SavedStateHandle is backed by an in-memory map
    // in unit tests.
    private data object InvalidSecondDestinationNavKey: NavKey

    @Test
    fun navigateTo_withInvalidSecondDestination_throwsIllegalArgumentException() = runTest {
        // Given a fresh NavigatorImpl whose back stack contains only ValidStartDestinationNavKey
        val navigator = NavigatorImpl(SavedStateHandle(), ValidStartDestinationNavKey)
        advanceUntilIdle()

        // When navigateTo is called with InvalidSecondDestinationNavKey
        assertThrows(IllegalArgumentException::class.java) {
            navigator.navigateTo(InvalidSecondDestinationNavKey)
        }
        advanceUntilIdle()

        // Then the back stack contains ValidStartDestinationNavKey
        val backStack = navigator.getBackStack()
        assertEquals(1, backStack.size)
        assertEquals(ValidStartDestinationNavKey, backStack[0])
    }
}
