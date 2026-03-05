package xyz.fjarm.navigationlib

import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.snapshotFlow
import androidx.compose.runtime.snapshots.SnapshotStateList
import androidx.compose.runtime.toMutableStateList
import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.navigation3.runtime.NavKey
import dagger.hilt.android.lifecycle.ActivityRetainedSavedState
import dagger.hilt.android.scopes.ActivityRetainedScoped
import kotlinx.coroutines.flow.launchIn
import kotlinx.coroutines.flow.onEach
import xyz.fjarm.loginandsignupfeatlib.LoginAndSignUpNavKey
import javax.inject.Inject

@ActivityRetainedScoped
class NavigatorImpl @Inject constructor(
    @ActivityRetainedSavedState private val savedStateHandle: SavedStateHandle,
): ViewModel(), Navigator {

    companion object {
        private const val KEY_BACKSTACK = "backstack"
    }

    private val _backStack: SnapshotStateList<NavKey> = savedStateHandle
        .get<ArrayList<NavKey>>(KEY_BACKSTACK)
        ?.toMutableStateList()
        ?: mutableStateListOf(LoginAndSignUpNavKey)

    init {
        // Side effect that observes changes to _backStack and saves to SavedStateHandle.
        snapshotFlow { _backStack.toList() }
            .onEach {
                savedStateHandle[KEY_BACKSTACK] = ArrayList(it)
            }
            .launchIn(viewModelScope)
    }

    override fun back(): Boolean {
        val last = _backStack.removeLastOrNull()
        return last != null
    }

    override fun clear(): Boolean {
        _backStack.clear()
        // LoginAndSignUpNavKey is the default start destination.
        return _backStack.add(LoginAndSignUpNavKey)
    }

    override fun getBackStack(): List<NavKey> {
        return _backStack.toList()
    }

    override fun navigateTo(destination: NavKey) {
        _backStack.add(destination)
    }
}
