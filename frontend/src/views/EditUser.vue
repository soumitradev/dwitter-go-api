<template>
  <div class="editUser">
    <p v-if="loading">Loading...</p>
    <EditUser v-else-if="result" v-bind="result.user" viewUser="randomGuy" />
    <p v-if="error">{{ error }}</p>
  </div>
</template>


<script>
// @ is an alias to /src
import User from "../components/User.vue";
import { useQuery } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import { useRoute } from 'vue-router'

export default {
  name: "EditUserView",
  setup() {
    const route = useRoute()

    const { result, loading, error } = useQuery(gql`
      query getUser($username: String!){
        user(username: $username) {
          username
          name
          email
          bio
          pfpURL
        }
      }
    `, {
      username: route.params.id
    })

    return {
      result,
      loading,
      error,
    }
  },
  created() {
  },
  components: {
    User,
  },
  methods: {
    login: function () {

    },
  }
};
</script>
