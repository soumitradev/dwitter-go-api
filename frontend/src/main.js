import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

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

createApp(App).use(router).mount('#app')
