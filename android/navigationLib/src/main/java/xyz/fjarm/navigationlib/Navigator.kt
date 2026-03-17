package xyz.fjarm.navigationlib

import androidx.navigation3.runtime.NavKey

interface Navigator {

    fun processSideEffect(sideEffect: NavigatorSideEffect)

    /**
     * Get the current back stack.
     *
     * @return The current back stack.
     */
    fun getBackStack(): List<NavKey>
}
