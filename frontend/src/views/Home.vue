<template>
  <div class="home">
    <!-- Nested objects are a bit weird rn, but we'll see -->
    <User
      v-bind="{
        bio: 'OH MY GOD I LOVE COCK OH GOD OH FUCK AAAAAAAAAAA AAAAAAA',
        createdAt: '2021-12-04T10:46:11.265Z',
        dweets: [
          {
            author: {
              bio: '',
              createdAt: '2021-12-04T10:46:11.265Z',
              email: 'ok@ok.com',
              followerCount: 1,
              followingCount: 0,
              name: 'OOOOK',
              pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
              username: 'OkayMan'
            },
            authorID: 'OkayMan',
            dweetBody: 'duderinos',
            id: 'sREiME3iw6',
            isReply: false,
            lastUpdatedAt: '2021-12-04T10:46:47.393Z',
            likeCount: 0,
            media: [],
            originalReplyID: '',
            postedAt: '2021-12-04T10:46:47.393Z',
            redweetCount: 1,
            replyCount: 0
          },
          {
            author: {
              bio: '',
              createdAt: '2021-12-04T10:46:11.265Z',
              email: 'ok@ok.com',
              followerCount: 1,
              followingCount: 0,
              name: 'OOOOK',
              pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
              username: 'OkayMan'
            },
            authorID: 'OkayMan',
            dweetBody: 'my',
            id: 'SYqNUB1J0r',
            isReply: false,
            lastUpdatedAt: '2021-12-04T10:46:41.521Z',
            likeCount: 0,
            media: [],
            originalReplyID: '',
            postedAt: '2021-12-04T10:46:41.521Z',
            redweetCount: 0,
            replyCount: 0
          },
          {
            author: {
              bio: '',
              createdAt: '2021-12-04T10:46:11.265Z',
              email: 'ok@ok.com',
              followerCount: 1,
              followingCount: 0,
              name: 'OOOOK',
              pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
              username: 'OkayMan'
            },
            authorID: 'OkayMan',
            dweetBody: 'epic',
            id: '2eUnRZ9XJN',
            isReply: false,
            lastUpdatedAt: '2021-12-04T10:46:37.767Z',
            likeCount: 0,
            media: [],
            originalReplyID: '',
            postedAt: '2021-12-04T10:46:37.767Z',
            redweetCount: 0,
            replyCount: 0
          },
          {
            author: {
              bio: '',
              createdAt: '2021-12-04T10:46:11.265Z',
              email: 'ok@ok.com',
              followerCount: 1,
              followingCount: 0,
              name: 'OOOOK',
              pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
              username: 'OkayMan'
            },
            authorID: 'OkayMan',
            dweetBody: 'is',
            id: 'DrWVkyzWsJ',
            isReply: false,
            lastUpdatedAt: '2021-12-04T10:46:34.228Z',
            likeCount: 0,
            media: [],
            originalReplyID: '',
            postedAt: '2021-12-04T10:46:34.228Z',
            redweetCount: 0,
            replyCount: 0
          },
          {
            author: {
              bio: '',
              createdAt: '2021-12-04T10:46:11.265Z',
              email: 'ok@ok.com',
              followerCount: 1,
              followingCount: 0,
              name: 'OOOOK',
              pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
              username: 'OkayMan'
            },
            authorID: 'OkayMan',
            dweetBody: 'cock',
            id: 'yS1IKn5xbA',
            isReply: false,
            lastUpdatedAt: '2021-12-04T10:46:30.748Z',
            likeCount: 0,
            media: [],
            originalReplyID: '',
            postedAt: '2021-12-04T10:46:30.748Z',
            redweetCount: 0,
            replyCount: 0
          }
        ],
        email: '',
        feedObjects: [],
        followerCount: 1,
        followers: [
          {
            bio: '',
            createdAt: '2021-12-04T10:47:38.95Z',
            email: 'okman@ok.com',
            followerCount: 0,
            followingCount: 1,
            name: 'OOOOK',
            pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayManALT'
          },
      
        ],
        following: [],
        followingCount: 0,
        likedDweets: [],
        name: 'OOOOK',
        pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
        redweetedDweets: [],
        redweets: [],
        username: 'OkayMan'
      }"
    />
  </div>
</template>

<script>
// @ is an alias to /src
import User from "../components/User.vue";

export default {
  name: "Home",
  components: {
    User,
  },
};
</script>
