package xyz.fjarm.navigation

import androidx.navigation3.runtime.NavKey

interface Navigation {

    fun processSideEffect(sideEffect: NavigationSideEffect)

    /**
     * Get the current back stack.
     *
     * @return The current back stack.
     */
    fun getBackStack(): List<NavKey>
}
