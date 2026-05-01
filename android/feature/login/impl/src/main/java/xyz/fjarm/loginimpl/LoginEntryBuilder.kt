package xyz.fjarm.loginimpl

import androidx.navigation3.runtime.EntryProviderScope
import androidx.navigation3.runtime.NavKey
import xyz.fjarm.loginapi.LoginNavKey
import xyz.fjarm.navigation.Navigation


fun EntryProviderScope<NavKey>.loginEntryBuilder(navigation: Navigation) {
    entry<LoginNavKey> {
        LoginScreen()
    }
}