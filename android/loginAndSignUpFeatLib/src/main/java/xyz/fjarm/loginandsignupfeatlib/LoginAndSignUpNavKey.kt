package xyz.fjarm.loginandsignupfeatlib

import android.os.Parcelable
import androidx.navigation3.runtime.NavKey
import kotlinx.parcelize.Parcelize
import kotlinx.serialization.Serializable

@Parcelize // Generate a Parcelable implementation for the object.
@Serializable // If using rememberNavBackStack, ensure that the object is serializable.
data object LoginAndSignUpNavKey : NavKey, Parcelable
