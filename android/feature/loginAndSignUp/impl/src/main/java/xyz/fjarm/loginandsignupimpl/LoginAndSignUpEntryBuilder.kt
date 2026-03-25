package xyz.fjarm.loginandsignupimpl

import androidx.navigation3.runtime.EntryProviderScope
import androidx.navigation3.runtime.NavKey
import xyz.fjarm.loginandsignupapi.LoginAndSignUpNavKey
import xyz.fjarm.loginapi.LoginNavKey
import xyz.fjarm.navigation.Navigation
import xyz.fjarm.navigation.NavigationSideEffect

fun EntryProviderScope<NavKey>.loginAndSignUpEntryBuilder(navigation: Navigation) {
    entry<LoginAndSignUpNavKey> {
        LoginAndSignUpScreen(
            navigateToLogIn = { navigation.processSideEffect(
                NavigationSideEffect.NavigateToDestination(LoginNavKey))
            },
            navigateToSignUp = {
                // TODO(2026-03-24): navigate to sign up screen
            }
        )
    }
}
