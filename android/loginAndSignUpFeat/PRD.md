## Feature summary
The log in and sign up feature enables returning Fjarm users to navigate to login and new Fjarm
users to navigate to sign up.

It's the first entry point into the Fjarm app and is shown when a user isn't logged in (no session
token is saved).

## User stories

As an unauthenticated Fjarm user, I want to navigate to login if I have existing credentials and
sign up if I'm a new user, so that I can enter the app and use its features.

## UI overview

| Log in and sign up landing screen                                             |
|-------------------------------------------------------------------------------|
| [![Log in and sign up landing](login-and-sign-up.png)](login-and-sign-up.png) |

## High level architecture - client

**Landing screen**

State:
* Title line
* Subtitle line
* Logo
* Navigate to sign up button text
* Navigate to log in button text

Events:
* Click sign up button
* Click log in button

Side effects:
* Navigate to sign up screen
* Navigate to log in screen

## High level architecture - API

No API requests need to be made by this screen or its `ViewModel`.

## Detailed architecture - client

`LoginAndSignUpFeatContract.kt`:

```kotlin
data class LoginAndSignUpState(
    @StringRes val titleLine: Int,
    @StringRes val subtitleLine: Int,
    @DrawableRes val logo: Int,
    @StringRes val signUpButtonText: Int,
    @StringRes val logInButtonText: Int,
)

sealed class LoginAndSignUpEvent {

    data object SignUpButtonClicked : LoginAndSignUpEvent()

    data object LogInButtonClicked : LoginAndSignUpEvent()
}

sealed class LoginAndSignUpSideEffect {
    data object NavigateToSignUp : LoginAndSignUpSideEffect()

    data object NavigateToLogIn : LoginAndSignUpSideEffect()
}
```

## Detailed architecture - API

No API requests need to be made by this screen or its `ViewModel`.

## Testing - client

- [ ] Given an unauthenticated user, when they navigate to the `:loginAndSignUpFeat` screen, then they see a title line, subtitle line, logo, sign up button, and log in button with correct text.
- [ ] Given an unauthenticated user, when they click the sign up button, then they navigate to the `:signUpFeat` screen.
- [ ] Given an unauthenticated user, when they click the log in button, then they navigate to the `:logInFeat` screen.

## Monitoring - client

**Screen performance**
- [ ] Time to first composition
- [ ] Recomposition counts (detect unnecessary recompositions)
- [ ] Frame render time and jank detection
- [ ] Screen load duration

**User actions**
- [ ] Log in button clicks
- [ ] Sign up button clicks
- [ ] Back button clicks

**Navigation flows**
- [ ] Entry point tracking
- [ ] Exit point tracking

## TODO

- [ ] Support log in and sign up with Google
- [ ] Support log in and sign up with Apple
- [ ] Include privacy policy
- [ ] Include ToS
