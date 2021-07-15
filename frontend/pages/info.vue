<template>
  <div>
    <section id="profile">
      browserLanguage : {{ browserLanguage }} <br>
      sdkVersion : {{ sdkVersion }} <br>
      isInClient : {{ isInClient }} <br>
      isLoggedIn : {{ isLoggedIn }} <br>
      deviceOS : {{ deviceOS }} <br>
      <img
        id="pictureUrl"
        :src="displayName"
        height="100"
      >
      <p id="userId">
        UserId: {{ userId }}
      </p>
      <p id="statusMessage">
        statusMessage: {{ statusMessage }}
      </p>
      <p id="email">
        email: {{ email }}
      </p>
      <button type="button" class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 md:py-4 md:text-lg md:px-10" @click="login()">
        Login
      </button>
    </section>
  </div>
</template>

<script>
// import liff from '@line/liff'

export default {
  data () {
    return {
      userId: '',
      displayName: '',
      statusMessage: '',
      email: '',
      accessToken: null,
      browserLanguage: null,
      sdkVersion: null,
      isInClient: null,
      isLoggedIn: null,
      deviceOS: null,
      l: null
    }
  },
  beforeCreate () {
  },
  mounted () {
    liff.init({ liffId: '1656205141-yPMrOQvn' })
      .then(() => {
        this.browserLanguage = liff.getLanguage()
        this.sdkVersion = liff.getVersion()
        this.isInClient = liff.isInClient()
        this.isLoggedIn = liff.isLoggedIn()
        this.deviceOS = liff.getOS()
        if (!liff.isLoggedIn()) {
          liff.login()
        }
      })
      .catch((err) => {
      // Error happens during initialization
      // console.log(err.code, err.message);
        alert(err.message)
      })
    if (liff.isLoggedIn()) {
      const profile = liff.getProfile()
      this.userId = profile.userId
      this.displayName = profile.displayName
      this.statusMessage = profile.statusMessage
    }
    // this.email = liff.getDecodedIDToken().email
  },
  methods: {
    login () {
      //  liff.login()
      liff.login({
        redirectUri: 'https://ituyen.herokuapp.com/info'
      })
    }
  }
}
</script>
