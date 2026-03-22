package xyz.fjarm.loginandsignupimpl

import androidx.navigation3.runtime.EntryProviderScope
import androidx.navigation3.runtime.NavKey
import xyz.fjarm.loginandsignupapi.LoginAndSignUpNavKey

fun EntryProviderScope<NavKey>.loginAndSignUpEntryBuilder() {
    entry<LoginAndSignUpNavKey> {
        LoginAndSignUpScreen()
    }
}
