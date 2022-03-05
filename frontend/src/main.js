import { createApp, provide, h } from 'vue'
import App from './App.vue'
import './styles/app.css'
import router from './router'
import { DefaultApolloClient } from '@vue/apollo-composable'
import { ApolloClient, ApolloLink, createHttpLink, InMemoryCache } from '@apollo/client/core'
import { setContext } from '@apollo/client/link/context';
import { TokenRefreshLink } from "apollo-link-token-refresh"
import decodeJWT from "jwt-decode"




// HTTP connection to the API
const httpLink = createHttpLink({
  uri: 'http://localhost:5000/api/graphql',
})

// Cache implementation
const cache = new InMemoryCache()

// Use auth
const authLink = setContext((_, { headers }) => {
  const token = localStorage.getItem('token');
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : "",
    }
  }
});

const refreshTokenLink = new TokenRefreshLink({
  isTokenValidOrUndefined: () => {
    const token = localStorage.getItem("token");
    if (!token) return true;

    const claims = decodeJWT(token)
    const expirationTimeInSeconds = claims.exp * 1000
    const now = new Date()
    return (expirationTimeInSeconds >= now.getTime())
  },
  fetchAccessToken: async () => {
    const refreshToken = localStorage.getItem("jid")
    const request = await fetch("http://localhost:5000/api/refresh_token", {
      method: "POST",
      credentials: "omit",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        jid: refreshToken,
      }),
    })
    const result = await request.json()
    return result
  },
  handleFetch: (accessToken) => {
    localStorage.setItem("token", accessToken)
  },
  handleResponse: (operation, accessTokenField) => (response) => {
    console.log(operation, accessTokenField, response)
    return { access_token: response.accessToken }
  },
  handleError: (err) => {
    console.warn("Error refreshing token.")
    console.error(err)
    localStorage.removeItem("jwt")
    localStorage.removeItem("refreshToken")
  },
})

// Create the apollo client
const apolloClient = new ApolloClient({
  link: ApolloLink.from([refreshTokenLink, authLink, httpLink]),
  cache,
})

const app = createApp({
  setup() {
    provide(DefaultApolloClient, apolloClient)
  },

  render: () => h(App),
})

/*
- 4th March: Finish basic frontend-backend connection
- 5th March: Finish pagination
- 6th March: Sign Up page
- 7th March: Finish EditUser connection
- 8th March: Figure out where to do email verification from, and maybe even implement it
- 9th March: slep
- 10th March: Fork dwitter meme and add learn basic redis shit
- 11th March: Move to session auth
- 12th March: 
- 13th March

*/

/*
Components to make

Utility
- Full Frame for the website Frame.vue
- Confirmation modal dialogue for 'delete dweet' or similar actions Modal.vue
- Some kind of settings view Settings.vue

Unauthenticated
- Sign up view SignUp.vue

Unauthenticated + Authenticated
- Dweet/Redweet in dweet View Dweet.vue
- Replies in dweet view Reply.vue
- User view User.vue

Only Authenticated
- Replies in feed view FeedReply.vue
- Dweet/Redweet in feed view FeedDweet.vue
- Writing reply view NewReply.vue
- Writing dweet view NewDweet.vue
- User list view (for followers/following) UserList.vue
- Edit user view EditUser.vue
- Edit dweet view EditDweet.vue
*/

app.use(router).mount('#app')
