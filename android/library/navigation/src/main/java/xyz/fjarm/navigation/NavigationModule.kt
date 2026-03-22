package xyz.fjarm.navigation

import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.android.components.ActivityRetainedComponent
import dagger.hilt.android.scopes.ActivityRetainedScoped

@Module
@InstallIn(ActivityRetainedComponent::class)
interface NavigationModule {

    @Binds
    @ActivityRetainedScoped
    fun bindsNavigation(impl: NavigationImpl): Navigation
}
