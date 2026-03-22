package xyz.fjarm.loginandsignupimpl

import androidx.navigation3.runtime.EntryProviderScope
import androidx.navigation3.runtime.NavKey
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.components.ActivityRetainedComponent
import dagger.multibindings.IntoSet

@Module
@InstallIn(ActivityRetainedComponent::class)
object LoginAndSignUpModule {

    @Provides
    @IntoSet
    fun provideLoginAndSignUpEntryBuilder(): EntryProviderScope<NavKey>.() -> Unit = {
        loginAndSignUpEntryBuilder()
    }
}
