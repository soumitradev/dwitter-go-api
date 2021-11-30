<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png" />
    <Dweet dweetBody="Welcome to Your Vue.js App" />
  </div>
</template>

<script>
// @ is an alias to /src
import Dweet from "../components/Dweet.vue";

export default {
  name: "Home",
  components: {
    Dweet,
  },
};
</script>
