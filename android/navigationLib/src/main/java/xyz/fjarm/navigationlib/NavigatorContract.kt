package xyz.fjarm.navigationlib

import androidx.navigation3.runtime.NavKey

sealed class NavigatorSideEffect {

    /**
     * Clear the back stack and navigate to the start destination.
     */
    data object ClearBackStack: NavigatorSideEffect()

    /**
     * Navigate back to the last destination in the back stack.
     */
    data object NavigateBack: NavigatorSideEffect()

    /**
     * Navigate to a specific destination.
     *
     * @param destination The destination to navigate to.
     */
    data class NavigateToDestination(val destination: NavKey): NavigatorSideEffect()
}
