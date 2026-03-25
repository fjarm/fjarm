package xyz.fjarm.loginandsignupimpl

import androidx.navigation3.runtime.EntryProviderScope
import androidx.navigation3.runtime.NavKey
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.components.ActivityRetainedComponent
import dagger.multibindings.IntoSet
import xyz.fjarm.navigation.Navigation

@Module
@InstallIn(ActivityRetainedComponent::class)
object LoginAndSignUpModule {

    @Provides
    @IntoSet
    fun provideLoginAndSignUpEntryBuilder(
        navigation: Navigation,
    ): EntryProviderScope<NavKey>.() -> Unit = {
        loginAndSignUpEntryBuilder(navigation)
    }
}
