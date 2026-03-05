package xyz.fjarm.navigationlib

import androidx.navigation3.runtime.NavKey

interface Navigator {

    /**
     * Navigate back to the last destination in the back stack.
     *
     * @return `true` if a destination was popped, `false` otherwise.
     */
    fun back(): Boolean

    /**
     * Clear the back stack and navigate to the start destination.
     *
     * @return `true` if the start destination was reached, `false` otherwise.
     */
    fun clear(): Boolean

    /**
     * Get the current back stack.
     *
     * @return The current back stack.
     */
    fun getBackStack(): List<NavKey>

    /**
     * Navigate to a specific destination.
     *
     * @param destination The destination to navigate to.
     */
    fun navigateTo(destination: NavKey)
}
