<template>
  <div class="about">
    <!-- <User /> -->
    <p v-if="error">{{ error }}</p>
    <p v-else-if="result">{{ result.user.username }}</p>
    <p v-else>Loading...</p>
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
  name: "User",
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
