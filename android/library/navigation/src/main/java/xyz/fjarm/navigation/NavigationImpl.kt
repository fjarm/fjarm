package xyz.fjarm.navigation

import android.os.Parcelable
import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.snapshots.SnapshotStateList
import androidx.compose.runtime.toMutableStateList
import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.navigation3.runtime.NavKey
import dagger.hilt.android.lifecycle.ActivityRetainedSavedState
import dagger.hilt.android.scopes.ActivityRetainedScoped
import javax.inject.Inject

@ActivityRetainedScoped
class NavigationImpl @Inject constructor(
    @ActivityRetainedSavedState private val savedStateHandle: SavedStateHandle,
    @StartDestination private val startDestination: NavKey,
): ViewModel(), Navigation {

    companion object {
        private const val KEY_BACKSTACK = "backstack"
    }

    private val _backStack: SnapshotStateList<NavKey> = savedStateHandle
        .get<ArrayList<NavKey>>(KEY_BACKSTACK)
        ?.toMutableStateList()
        ?: mutableStateListOf(startDestination)

    init {
        persist()
    }

    override fun processSideEffect(sideEffect: NavigationSideEffect) {
        when (sideEffect) {
            is NavigationSideEffect.ClearBackStack -> {
                _backStack.clear()
                // LoginAndSignUpNavKey is the default start destination.
                val success = _backStack.add(startDestination)
                if (success) persist()
            }
            is NavigationSideEffect.NavigateBack -> {
                val last = _backStack.removeLastOrNull()
                if (last != null) persist()
            }
            is NavigationSideEffect.NavigateToDestination -> {
                require(sideEffect.destination is Parcelable) {
                    "NavKey ${sideEffect.destination::class.simpleName} must implement Parcelable."
                }
                _backStack.add(sideEffect.destination)
                persist()
            }
        }
    }

    override fun getBackStack(): List<NavKey> {
        return _backStack.toList()
    }

    private fun persist() {
        savedStateHandle[KEY_BACKSTACK] = ArrayList(_backStack)
    }
}
