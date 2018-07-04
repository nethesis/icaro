import Auth0Lock from "auth0-lock";

var LoginService = {
  data() {
    return {
      auth0Lock: null
    }
  },
  methods: {
    initAuth0Lock() {
      this.auth0Lock = new Auth0Lock(
        CONFIG.AUTH0_CLIENT_ID,
        CONFIG.AUTH0_DOMAIN, {
          container: "auth0-login",
          socialButtonStyle: 'small',
          allowedConnections: [
            "github",
            "google-oauth2",
            "linkedin",
            "windowslive"
          ],
          auth: {
            responseType: "token id_token",
            params: {
              scope: "openid email profile"
            },
            audience: CONFIG.AUTH0_AUDIENCE,
            autoParseHash: true,
            redirect: true
          }
        }
      );

      return this.auth0Lock;
    },
    showAuth0Lock() {
      this.auth0Lock.show();
    },
    handleAuth0Lock(accessToken) {
      var context = this
      this.auth0Lock.getUserInfo(accessToken, function (error, profile) {
        if (error) {
          // Handle error
          return;
        }

        // save profile data
        context.set("auth0User", profile)

        // login with auth0
        context.execLoginAuth0(
          accessToken,
          success => {
            // extract loggedUser info
            var loggedUser = success.body;
            // save to localstorage
            context.set("loggedUser", loggedUser);

            // get user info
            context.getInfo(loggedUser.id, response => {
              if (response) {
                context.user.info = response;
                context.isLogged = true;
                context.initGraphics();
                context.user.info.name = profile.name;
                context.user.info.email = profile.email;
                context.user.info.picture = profile.picture;

                // change route
                context.isLogged = true;
                context.delete("auth0Data");
                context.$router.replace({
                  path: "/"
                })
              } else {
                context.isLogged = false;
              }
            });
          },
          error => {
            console.log(error.body.message);
          }
        );
      });
    },
    execLogin(body, success, error) {
      this.$http
        .post(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/login",
          body
        )
        .then(success, error);
    },
    execLoginAuth0(access_token, success, error) {
      this.$http
        .post(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/login-auth0", {}, {
            headers: {
              Authorization: "Bearer " + access_token || ""
            }
          }
        )
        .then(success, error);
    },
    execLogout(success, error) {
      this.$http
        .post(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/logout", {}, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    execGetInfo(id, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/accounts/" +
          id, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    execChangePassword(password, id, success, error) {
      this.$http
        .put(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/accounts/" +
          id, {
            password: password
          }, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    }
  }
};
export default LoginService;
