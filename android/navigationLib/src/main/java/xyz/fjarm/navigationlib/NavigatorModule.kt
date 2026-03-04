package xyz.fjarm.navigationlib

import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.components.ActivityRetainedComponent
import dagger.hilt.android.scopes.ActivityRetainedScoped
import xyz.fjarm.loginandsignupfeatlib.LoginAndSignUpNavKey

@Module
@InstallIn(ActivityRetainedComponent::class)
class NavigatorModule {

    @Provides
    @ActivityRetainedScoped
    fun provideNavigator(): Navigator {
        return NavigatorImpl(LoginAndSignUpNavKey)
    }
}
