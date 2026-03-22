package xyz.fjarm.fjarm

import androidx.navigation3.runtime.NavKey
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.components.ActivityRetainedComponent
import xyz.fjarm.loginandsignupapi.LoginAndSignUpNavKey
import xyz.fjarm.navigation.StartDestination

@Module
@InstallIn(ActivityRetainedComponent::class)
object FjarmModule {

    @Provides
    @StartDestination
    fun provideStartDestination(): NavKey = LoginAndSignUpNavKey
}
