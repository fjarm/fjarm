package xyz.fjarm.navigationlib

import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.snapshots.SnapshotStateList
import androidx.navigation3.runtime.NavKey
import dagger.hilt.android.scopes.ActivityRetainedScoped

@ActivityRetainedScoped
class NavigatorImpl(
    private val startDestination: NavKey,
): Navigator {

    private val backStack: SnapshotStateList<NavKey> = mutableStateListOf(startDestination)

    override fun back(): Boolean {
        val last = backStack.removeLastOrNull()
        return last != null
    }

    override fun clear(): Boolean {
        backStack.clear()
        return backStack.add(startDestination)
    }

    override fun getBackStack(): SnapshotStateList<NavKey> {
        return backStack
    }

    override fun navigateTo(destination: NavKey) {
        backStack.add(destination)
    }
}
