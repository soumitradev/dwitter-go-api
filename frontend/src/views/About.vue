<template>
  <div class="about">
    <NewDweet
      :viewUser="{
        bio: '',
        createdAt: '2021-12-04T10:46:11.265Z',
        email: 'ok@ok.com',
        followerCount: 1,
        followingCount: 0,
        name: 'OOOOK',
        pfpURL: 'https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media',
        username: 'OkayMan'
      }"
    />
  </div>
</template>


<script>
// @ is an alias to /src
import NewDweet from "../components/NewDweet.vue";

export default {
  name: "About",
  components: {
    NewDweet,
  },
};
</script>
