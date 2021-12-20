<template>
  <div class="home">
    <!-- Nested objects are a bit weird rn, but we'll see -->
    <Dweet
      v-bind="{
        author: {
          bio: '',
          createdAt: '2021-12-04T10:46:11.265Z',
          email: 'ok@ok.com',
          followerCount: 1,
          followingCount: 0,
          name: 'OOOOK',
          pfpURL:
            'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
          username: 'OkayMan',
        },
        authorID: 'OkayMan',
        dweetBody: 'duderinos',
        id: 'sREiME3iw6',
        isReply: false,
        lastUpdatedAt: '2021-12-04T10:46:47.393Z',
        likeCount: 0,
        likeUsers: [
          {
            bio: '',
            createdAt: '2021-12-04T10:46:11.265Z',
            email: 'ok@ok.com',
            followerCount: 1,
            followingCount: 0,
            name: 'OOOOK',
            pfpURL:
              'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayMan',
          },
        ],
        media: [],
        originalReplyID: '',
        postedAt: '2021-12-04T10:46:47.393Z',
        redweetCount: 1,
        redweetUsers: [
          {
            bio: '',
            createdAt: '2021-12-04T10:46:11.265Z',
            email: 'ok@ok.com',
            followerCount: 1,
            followingCount: 0,
            name: 'RedweetManPagMan',
            pfpURL:
              'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
            username: 'OkayMan',
          },
        ],
        replyCount: 0,
        replyDweets: [],
        replyTo: {
          author: {
            bio: '',
            createdAt: '0001-01-01T00:00:00Z',
            email: '',
            followerCount: 0,
            followingCount: 0,
            name: '',
            pfpURL: '',
            username: '',
          },
          authorID: '',
          dweetBody: '',
          id: '',
          isReply: false,
          lastUpdatedAt: '0001-01-01T00:00:00Z',
          likeCount: 0,
          media: [],
          originalReplyID: '',
          postedAt: '0001-01-01T00:00:00Z',
          redweetCount: 0,
          replyCount: 0,
        },
      }"
    />
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
