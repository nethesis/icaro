<template>
  <div class="ui segment">
    <div v-if="!hotspot.onError">
      <h3>Welcome</h3>
      <img class="ui centered small image" src="../assets/logo.png">
      <p>
        Quo homero blandit intellegebat. Incorrupte consequuntur mei id. Mei ut facer dolores adolescens, no illum aperiri quo, usu
        odio brute at
      </p>
      <div class="ui divider"></div>
      <router-link to='/login' class="big ui green button">
        Start Navigate
      </router-link>
    </div>
    <div v-if="hotspot.onError" class="ui icon negative message">
      <i class="warning icon"></i>
      <div class="content">
        <div class="header">
          Some params are missing. Can't we go ahead
        </div>
        <p>Something gone wrong opening the page...</p>
      </div>
    </div>
  </div>

</template>

<script>
  export default {
    name: 'SplashPage',
    data() {
      return {
        hotspot: {
          onError: this.$root.$options.hotspot.onError
        }
      }
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  h1,
  h2 {
    font-weight: normal;
  }

  ul {
    list-style-type: none;
    padding: 0;
  }

  li {
    display: inline-block;
    margin: 0 10px;
  }

  a {
    color: #42b983;
  }
</style>