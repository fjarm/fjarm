package xyz.fjarm.navigationlib

import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.android.components.ActivityRetainedComponent
import dagger.hilt.android.scopes.ActivityRetainedScoped

@Module
@InstallIn(ActivityRetainedComponent::class)
abstract class NavigatorModule {

    @Binds
    @ActivityRetainedScoped
    abstract fun bindsNavigator(impl: NavigatorImpl): Navigator
}
