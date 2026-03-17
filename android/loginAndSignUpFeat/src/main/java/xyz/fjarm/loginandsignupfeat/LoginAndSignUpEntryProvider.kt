package xyz.fjarm.loginandsignupfeat

import androidx.navigation3.runtime.EntryProviderScope
import androidx.navigation3.runtime.NavKey
import xyz.fjarm.loginandsignupfeatlib.LoginAndSignUpNavKey

fun EntryProviderScope<NavKey>.loginAndSignUpEntryBuilder() {
    entry<LoginAndSignUpNavKey> {
        LoginAndSignUpScreen()
    }
}
