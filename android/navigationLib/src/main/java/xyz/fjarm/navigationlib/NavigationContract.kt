package xyz.fjarm.navigationlib

import androidx.navigation3.runtime.NavKey

sealed class NavigationSideEffect {

    /**
     * Clear the back stack and navigate to the start destination.
     */
    data object ClearBackStack: NavigationSideEffect()

    /**
     * Navigate back to the last destination in the back stack.
     */
    data object NavigateBack: NavigationSideEffect()

    /**
     * Navigate to a specific destination.
     *
     * @param destination The destination to navigate to.
     */
    data class NavigateToDestination(val destination: NavKey): NavigationSideEffect()
}
