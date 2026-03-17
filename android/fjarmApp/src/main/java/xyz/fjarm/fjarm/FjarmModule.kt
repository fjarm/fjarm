package xyz.fjarm.fjarm

import androidx.navigation3.runtime.NavKey
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.components.ActivityRetainedComponent
import dagger.hilt.android.scopes.ActivityRetainedScoped
import xyz.fjarm.loginandsignupfeatlib.LoginAndSignUpNavKey
import xyz.fjarm.navigationlib.StartDestination

@Module
@InstallIn(ActivityRetainedComponent::class)
class FjarmModule {

    @Provides
    @ActivityRetainedScoped
    @StartDestination
    fun provideStartDestination(): NavKey = LoginAndSignUpNavKey
}
