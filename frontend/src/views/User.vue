<template>
  <div class="about">
    <p v-if="loading">Loading...</p>
    <User v-else-if="result" v-bind="result.user" viewUser="randomGuy" />
    <p v-if="error">{{ error }}</p>
  </div>
</template>


<script>
// @ is an alias to /src
import User from "../components/User.vue";
import { useQuery } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import { useRoute } from 'vue-router'
import { userFrag } from "../fragments/userFrag";

export default {
  name: "ViewUser",
  setup() {
    const route = useRoute()

    const { result, loading, error } = useQuery(gql`
      query getUser($username: String!){
        user(username: $username feedObjectsToFetch: 5 objectsToFetch: "feed") {
          ...UserFrag
        }
      }
      ${userFrag}
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
  }
};
</script>
