<template>
  <div class="login">
    <div class="flex flex-col">
      <input
        class="bg-neutral-99 mt-2 mx-16 text-left text-neutral-10 border-neutralVariant-50 border-opacity-50 border-x-0 border-t-0 border-b-2 resize-none appearance-none outline-none"
        v-model.trim="username"
        autofocus
      />
      <input
        class="bg-neutral-99 mt-2 mx-16 text-left text-neutral-10 border-neutralVariant-50 border-opacity-50 border-x-0 border-t-0 border-b-2 resize-none appearance-none outline-none"
        v-model.trim="password"
        type="password"
      />
      <button class="bg-primary-70 mx-16 rounded-full m-4" @click="login()">Login</button>
    </div>
  </div>
</template>


<script>
// @ is an alias to /src
import User from "../components/User.vue";

export default {
  name: "EditUserView",
  setup() {
  },
  created() {
  },
  components: {
    User,
  },
  methods: {
    login: function () {
      fetch("http://localhost:5000/api/login", {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          username: this.username,
          password: this.password
        })
      }).then(res => res.json()).then(data => {
        localStorage.setItem('token', data.accessToken);
        localStorage.setItem('jid', data.jid);
      }).catch((error) => {
        console.error('Error:', error);
      });
    },
  }
};
</script>
