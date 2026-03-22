package xyz.fjarm.navigation

import android.os.Bundle
import android.os.Parcel
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
import kotlinx.serialization.Serializable
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertThrows
import org.junit.Before
import org.junit.Test
import org.junit.runner.RunWith
import org.robolectric.RobolectricTestRunner
import org.robolectric.annotation.Config

@OptIn(ExperimentalCoroutinesApi::class)
@RunWith(RobolectricTestRunner::class)
@Config(sdk = [31, 36])
class NavigationImplRobolectricTest {

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

    @Parcelize
    @Serializable
    private data object ValidSecondDestinationNavKey: NavKey, Parcelable

    @Serializable
    private data object InvalidSecondDestinationNavKey: NavKey

    @Test
    fun navigationImpl_withChangedStartDestination_backStackSurvivesProcessDeath() = runTest {
        // Given a fresh NavigationImpl whose back stack contains only ValidStartDestinationNavKey
        val savedStateHandle = SavedStateHandle()
        val navigation = NavigationImpl(savedStateHandle, ValidStartDestinationNavKey)
        advanceUntilIdle()

        val backStack = navigation.getBackStack()
        assertEquals(1, backStack.size)
        assertEquals(ValidStartDestinationNavKey, backStack[0])

        // Simulate process death: save state to a Bundle, then force a full Parcel round-trip. This
        // is the operation that enforces @Parcelize — any NavKey not annotated will throw here,
        // matching the production failure.
        val savedBundle = savedStateHandle.savedStateProvider().saveState()
        val parcel = Parcel.obtain()
        try {
            savedBundle.writeToParcel(parcel, 0)
            parcel.setDataPosition(0)
            val restoredBundle = Bundle.CREATOR.createFromParcel(parcel)

            // Restore a new NavigationImpl from the recovered Bundle
            val restoredHandle = SavedStateHandle.createHandle(restoredBundle, null)
            val restoredNavigation = NavigationImpl(restoredHandle, ValidSecondDestinationNavKey)
            advanceUntilIdle()

            // Then the restored back stack contains ValidStartDestinationNavKey
            val restoredBackStack = restoredNavigation.getBackStack()
            assertEquals(1, restoredBackStack.size)
            assertEquals(ValidStartDestinationNavKey, restoredBackStack[0])
        } catch (e: Exception) {
            throw AssertionError("Process death test failed: ${e.message}")
        } finally {
            parcel.recycle()
        }
    }

    @Test
    fun navigateTo_withValidSecondDestinationNavKey_doesNotThrowIllegalArgumentException() = runTest {
        // Given a fresh NavigationImpl whose back stack contains only ValidStartDestinationNavKey
        val savedStateHandle = SavedStateHandle()
        val navigation = NavigationImpl(savedStateHandle, ValidStartDestinationNavKey)
        advanceUntilIdle()

        // When navigateTo is called with ValidSecondDestinationNavKey which does implement Parcelable
        navigation.processSideEffect(NavigationSideEffect.NavigateToDestination(
            ValidSecondDestinationNavKey,
        ))
        advanceUntilIdle()

        // Simulate process death: save state to a Bundle, then force a full Parcel round-trip. This
        // is the operation that enforces @Parcelize — any NavKey not annotated will throw here,
        // matching the production failure.
        val savedBundle = savedStateHandle.savedStateProvider().saveState()
        val parcel = Parcel.obtain()
        try {
            savedBundle.writeToParcel(parcel, 0)
            parcel.setDataPosition(0)
            val restoredBundle = Bundle.CREATOR.createFromParcel(parcel)

            // Restore a new NavigationImpl from the recovered Bundle
            val restoredHandle = SavedStateHandle.createHandle(restoredBundle, null)
            val restoredNavigation = NavigationImpl(restoredHandle, ValidStartDestinationNavKey)
            advanceUntilIdle()

            // Then the restored back stack contains ValidStartDestinationNavKey and ValidDestinationNavKey
            val restoredBackStack = restoredNavigation.getBackStack()
            advanceUntilIdle()

            assertEquals(2, restoredBackStack.size)
            assertEquals(ValidStartDestinationNavKey, restoredBackStack[0])
            assertEquals(ValidSecondDestinationNavKey, restoredBackStack[1])
        } catch (e: Exception) {
            throw AssertionError("Process death test failed: ${e.message}")
        } finally {
            parcel.recycle()
        }
    }

    @Test
    fun navigateTo_withInvalidSecondDestinationNavKey_throwsIllegalArgumentException() = runTest {
        // Given a fresh NavigationImpl whose back stack contains only ValidStartDestinationNavKey
        val savedStateHandle = SavedStateHandle()
        val navigation = NavigationImpl(savedStateHandle, ValidStartDestinationNavKey)
        advanceUntilIdle()

        // When navigateTo is called with InvalidDestinationNavKey which does not implement
        // Parcelable
        // Then an IllegalArgumentException is thrown
        assertThrows(IllegalArgumentException::class.java) {
            navigation.processSideEffect(NavigationSideEffect.NavigateToDestination(
                InvalidSecondDestinationNavKey,
            ))
        }
    }
}
