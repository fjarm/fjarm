package xyz.fjarm.loginandsignupimpl

import androidx.navigation3.runtime.EntryProviderScope
import androidx.navigation3.runtime.NavKey
import xyz.fjarm.loginandsignupapi.LoginAndSignUpNavKey
import xyz.fjarm.navigation.Navigation
import xyz.fjarm.navigation.NavigationSideEffect

fun EntryProviderScope<NavKey>.loginAndSignUpEntryBuilder(navigation: Navigation) {
    entry<LoginAndSignUpNavKey> {
        LoginAndSignUpScreen(
            navigateToSignUp = { navigation.processSideEffect(
                NavigationSideEffect.NavigateToDestination())
            },
            navigateToLogIn = { navigation.processSideEffect(
                NavigationSideEffect.NavigateToDestination())
            }
        )
    }
}
