<template>
  <div class="editUser">
    <p v-if="loading">Loading...</p>
    <EditUser
      v-else-if="result"
      v-bind="result.user"
      :viewUser="parseJwt().username"
      @save-edits-to-user="saveEdits"
    />
    <p v-if="error">{{ error }}</p>
  </div>
</template>


<script>
// @ is an alias to /src
import EditUser from "../components/EditUser.vue";
import { useQuery } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import { userFrag } from "../fragments/userFrag";

export default {
  name: "EditUserView",
  setup() {
    let parseJwt = function () {
      let token = localStorage.getItem('token')
      var base64Url = token.split('.')[1];
      var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
      var jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
      }).join(''));
      return JSON.parse(jsonPayload);
    }

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
      username: parseJwt().username
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
    EditUser,
  },
  methods: {
    // https://stackoverflow.com/a/38935990
    urltoFile: async function (url, mimeType) {
      const res = await fetch(url);
      const data = await res.arrayBuffer();
      return new Blob([data], { type: mimeType });
    },
    parseJwt: function () {
      let token = localStorage.getItem('token')
      var base64Url = token.split('.')[1];
      var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
      var jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
      }).join(''));
      return JSON.parse(jsonPayload);
    },
    saveEdits: async function (value) {
      let updatedPfpURL = this.result.user.pfpURL;

      if (value.pfpUpdated) {
        const formData = new FormData();

        const imageData = await this.urltoFile(value.pfpURL, "image/png");
        formData.append('files', imageData);

        const token = localStorage.getItem("token")

        fetch("http://localhost:5000/api/pfp_upload", {
          method: 'POST',
          body: formData,
          headers: {
            authorization: token ? `Bearer ${token}` : "",
          }
        }).then(res => res.json()).then(data => {
          updatedPfpURL = data[0];

          this.$apollo.mutate({
            mutation: gql`mutation editUser($name: String!, $email: String!, $bio: String!, $pfpURL: String!){
              editUser(name: $name, email: $email, bio: $bio, pfpURL: $pfpURL) {
                ...UserFrag
              }
            }
            ${userFrag}`,
            variables: {
              name: value.name,
              email: value.email,
              bio: value.bio,
              pfpURL: updatedPfpURL,
            },
          }).then(
          );
        }).catch((error) => {
          console.error('Error:', error);
        });
        return;
      }

      this.$apollo.mutate({
        mutation: gql`mutation editUser($name: String!, $email: String!, $bio: String!, $pfpURL: String!){
          editUser(name: $name, email: $email, bio: $bio, pfpURL: $pfpURL) {
            ...UserFrag
          }
        }
        ${userFrag}`,
        variables: {
          name: value.name,
          email: value.email,
          bio: value.bio,
          pfpURL: updatedPfpURL,
        },
      }).then(
      );
    },
  }
};
</script>
